// Copyright (c) 2020-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package proxy

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	pluginapi "github.com/mattermost/mattermost-plugin-api"
	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/mattermost/mattermost-plugin-apps/server/api"
	"github.com/mattermost/mattermost-plugin-apps/server/api/impl/aws"
	"github.com/mattermost/mattermost-plugin-apps/server/api/impl/upstream"
	"github.com/mattermost/mattermost-plugin-apps/server/api/impl/upstream/upawslambda"
	"github.com/mattermost/mattermost-plugin-apps/server/api/impl/upstream/uphttp"
)

type Proxy struct {
	// Built-in Apps are linked in Go and invoked directly. The list is
	// initialized on startup, and need not be synchronized.
	builtinProvisionedApps map[api.AppID]api.Upstream

	mm        *pluginapi.Client
	conf      api.Configurator
	store     api.Store
	awsClient *aws.Client
}

var _ api.Proxy = (*Proxy)(nil)

func NewProxy(mm *pluginapi.Client, awsClient *aws.Client, conf api.Configurator, store api.Store) *Proxy {
	return &Proxy{
		mm:        mm,
		conf:      conf,
		store:     store,
		awsClient: awsClient,
	}
}

func (p *Proxy) Call(adminAccessToken string, call *api.Call) (*api.Call, *api.CallResponse) {
	fmt.Printf("<><> Call 1: %s %q\n", call.URL, call.Type)
	conf := p.conf.GetConfig()
	app, err := p.store.LoadApp(call.Context.AppID)
	if err != nil {
		return call, api.NewErrorCallResponse(err)
	}

	oauth := p.newMattermostOAuthenticator(app)
	call, err = p.expandCall(call, app, adminAccessToken, oauth, nil)
	if err == errOAuthRequired {
		connectURL := oauth.GetConnectURL()
		fmt.Printf("<><> Call 2: connectURL: %q, %v\n", connectURL, err)

		post := &model.Post{
			UserId:    conf.BotUserID,
			ChannelId: call.Context.ChannelID,
			Message:   fmt.Sprintf("If you are not automatically redirected, please click [here](%s) to connect.", connectURL),
		}
		p.mm.Post.SendEphemeralPost(call.Context.ActingUserID, post)
		err = p.mm.Post.DM(conf.BotUserID, call.Context.ActingUserID, post)
		fmt.Printf("<><> Call 3: %v\n", err)
		return call, &api.CallResponse{
			Type:          api.CallResponseTypeNavigate,
			NavigateToURL: connectURL,
		}
	}
	if err != nil {
		return call, api.NewErrorCallResponse(err)
	}

	up, err := p.upstreamForApp(app)
	if err != nil {
		return call, api.NewErrorCallResponse(err)
	}
	cr := upstream.Call(up, call)
	fmt.Printf("<><> Call 4: %s done: %q: %q %q\n", call.URL, cr.Type, cr.Markdown, cr.ErrorText)

	// TODO: the user-agents do not yet support Navigate, so post messages with the URL
	if cr.Type == api.CallResponseTypeNavigate {
		post := &model.Post{
			UserId:    conf.BotUserID,
			ChannelId: call.Context.ChannelID,
			Message:   fmt.Sprintf("If you are not automatically redirected, please navigate [here](%s) to continue.", cr.NavigateToURL),
		}
		p.mm.Post.SendEphemeralPost(call.Context.ActingUserID, post)
	}

	return call, cr
}

func (p *Proxy) Notify(cc *api.Context, subj api.Subject) error {
	subs, err := p.store.LoadSubs(subj, cc.TeamID, cc.ChannelID)
	if err != nil {
		return err
	}

	expandCache := &expandCache{}

	notify := func(sub *api.Subscription) error {
		call := sub.Call
		if call == nil {
			return errors.New("nothing to call")
		}
		app, err := p.store.LoadApp(sub.AppID)
		if err != nil {
			return err
		}
		oauth := p.newMattermostOAuthenticator(app)
		call, err = p.expandCall(call, app, "", oauth, expandCache)
		// TODO: DM the user to renew expired tokens?
		if err == errOAuthRequired {
			return errors.New("missing or invalid OAuth2 token")
		}
		if err != nil {
			return err
		}

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
		if len(p.builtinProvisionedApps) == 0 {
			return nil, errors.Errorf("builtin app not found: %s", app.Manifest.AppID)
		}
		up := p.builtinProvisionedApps[app.Manifest.AppID]
		if up == nil {
			return nil, errors.Errorf("builtin app not found: %s", app.Manifest.AppID)
		}
		return up, nil

	default:
		return nil, errors.Errorf("not a valid app type: %s", app.Manifest.Type)
	}
}

func (p *Proxy) ProvisionBuiltIn(appID api.AppID, up api.Upstream) {
	if p.builtinProvisionedApps == nil {
		p.builtinProvisionedApps = map[api.AppID]api.Upstream{}
	}
	p.builtinProvisionedApps[appID] = up
}

func WriteCallError(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(api.NewErrorCallResponse(err))
}
