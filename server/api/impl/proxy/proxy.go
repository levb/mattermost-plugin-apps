// Copyright (c) 2020-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package proxy

import (
	"encoding/json"
	"fmt"
	"net/http"

	pluginapi "github.com/mattermost/mattermost-plugin-api"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-apps/server/api"
	"github.com/mattermost/mattermost-plugin-apps/server/api/impl/aws"
	"github.com/mattermost/mattermost-plugin-apps/server/api/impl/upstream"
	"github.com/mattermost/mattermost-plugin-apps/server/api/impl/upstream/upawslambda"
	"github.com/mattermost/mattermost-plugin-apps/server/api/impl/upstream/uphttp"
)

type Proxy struct {
	builtIn map[api.AppID]api.Upstream

	mm        *pluginapi.Client
	conf      api.Configurator
	store     api.Store
	awsClient *aws.Client
}

var _ api.Proxy = (*Proxy)(nil)

func NewProxy(mm *pluginapi.Client, awsClient *aws.Client, conf api.Configurator, store api.Store) *Proxy {
	return &Proxy{nil, mm, conf, store, awsClient}
}

func (p *Proxy) Call(sessionToken api.SessionToken, call *api.Call) *api.CallResponse {
	conf := p.conf.GetConfig()
	app, err := p.store.LoadApp(call.Context.AppID)
	if err != nil {
		return api.NewErrorCallResponse(err)
	}

	expander := p.newExpander(call.Context, sessionToken)
	callContext, connectURL, err := p.ExpandForApp(expander, call, app)
	if err != nil {
		return api.NewErrorCallResponse(err)
	}
	if connectURL != "" {
		fmt.Printf("<><> Call 3: connectURL: %q, %v\n", connectURL, err)

		post := &model.Post{
			UserId:    conf.BotUserID,
			ChannelId: call.Context.ChannelID,
			Message:   fmt.Sprintf("If you are not automatically redirected, please click [here](%s) to connect.", connectURL),
		}
		p.mm.Post.SendEphemeralPost(call.Context.ActingUserID, post)
		err = p.mm.Post.DM(conf.BotUserID, call.Context.ActingUserID, post)
		fmt.Printf("<><> Call 4: %v\n", err)
		return &api.CallResponse{
			Type:          api.CallResponseTypeNavigate,
			NavigateToURL: connectURL,
		}
	}
	call.Context = callContext

	up, err := p.upstreamForApp(app)
	if err != nil {
		return api.NewErrorCallResponse(err)
	}
	cr := upstream.Call(up, call)
	fmt.Printf("<><> Call 2: cr: %+v\n", cr)

	// TODO: the user-agents do not yet support Navigate, so post messages with the URL
	if cr.Type == api.CallResponseTypeNavigate {
		post := &model.Post{
			UserId:    conf.BotUserID,
			ChannelId: call.Context.ChannelID,
			Message:   fmt.Sprintf("If you are not automatically redirected, please navigate [here](%s) to continue.", cr.NavigateToURL),
		}
		p.mm.Post.SendEphemeralPost(call.Context.ActingUserID, post)
	}

	return cr
}

func (p *Proxy) Notify(cc *api.Context, subj api.Subject) error {
	subs, err := p.store.LoadSubs(subj, cc.TeamID, cc.ChannelID)
	if err != nil {
		return err
	}

	expander := p.newExpander(cc, "")

	notify := func(sub *api.Subscription) error {
		call := sub.Call
		if call == nil {
			return errors.New("nothing to call")
		}
		app, err := p.store.LoadApp(sub.AppID)
		if err != nil {
			return err
		}
		callContext, connectURL, err := p.ExpandForApp(expander, call, app)
		if err != nil {
			return err
		}
		// TODO: DM the user to renew expired tokens?
		if connectURL != "" {
			return errors.New("missing or invalid OAuth2 token")
		}
		callContext.Subject = subj
		call.Context = callContext

		up, err := p.upstreamForApp(app)
		if err != nil {
			return err
		}
		return upstream.Notify(up, call)
	}

	for _, sub := range subs {
		err := notify(sub)
		if err != nil {
			// TODO log err
			continue
		}
	}
	return nil
}

func (p *Proxy) upstreamForApp(app *api.App) (api.Upstream, error) {
	switch app.Manifest.Type {
	case api.AppTypeHTTP:
		return uphttp.NewUpstream(app), nil

	case api.AppTypeAWSLambda:
		return upawslambda.NewUpstream(app, p.awsClient), nil

	case api.AppTypeBuiltin:
		if len(p.builtIn) == 0 {
			return nil, errors.Errorf("builtin app not found: %s", app.Manifest.AppID)
		}
		up := p.builtIn[app.Manifest.AppID]
		if up == nil {
			return nil, errors.Errorf("builtin app not found: %s", app.Manifest.AppID)
		}
		return up, nil

	default:
		return nil, errors.Errorf("not a valid app type: %s", app.Manifest.Type)
	}
}

func (p *Proxy) ProvisionBuiltIn(appID api.AppID, up api.Upstream) {
	if p.builtIn == nil {
		p.builtIn = map[api.AppID]api.Upstream{}
	}
	p.builtIn[appID] = up
}

func WriteCallError(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(api.NewErrorCallResponse(err))
}
