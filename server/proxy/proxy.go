// Copyright (c) 2020-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package proxy

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/server/upstream"
	"github.com/mattermost/mattermost-plugin-apps/server/upstream/upawslambda"
	"github.com/mattermost/mattermost-plugin-apps/server/upstream/uphttp"
)

func (p *proxy) Call(adminAccessToken string, call *apps.Call) (*apps.Call, *apps.CallResponse) {
	fmt.Printf("<><> Call 1: %s %q\n", call.URL, call.Type)
	conf := p.conf.Get()
	app, err := p.store.App.Get(call.Context.AppID)
	if err != nil {
		return call, apps.NewErrorCallResponse(err)
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
		return call, &apps.CallResponse{
			Type:          apps.CallResponseTypeNavigate,
			NavigateToURL: connectURL,
		}
	}
	if err != nil {
		return call, apps.NewErrorCallResponse(err)
	}

	up, err := p.upstreamForApp(app)
	if err != nil {
		return call, apps.NewErrorCallResponse(err)
	}
	cr := upstream.Call(up, call)
	fmt.Printf("<><> Call 4: %s done: %q: %q %q\n", call.URL, cr.Type, cr.Markdown, cr.ErrorText)

	// TODO: the user-agents do not yet support Navigate, so post messages with the URL
	if cr.Type == apps.CallResponseTypeNavigate {
		post := &model.Post{
			UserId:    conf.BotUserID,
			ChannelId: call.Context.ChannelID,
			Message:   fmt.Sprintf("If you are not automatically redirected, please navigate [here](%s) to continue.", cr.NavigateToURL),
		}
		p.mm.Post.SendEphemeralPost(call.Context.ActingUserID, post)
	}

	return call, cr
}

func (p *proxy) Notify(cc *apps.Context, subj apps.Subject) error {
	subs, err := p.store.Subscription.Get(subj, cc.TeamID, cc.ChannelID)
	if err != nil {
		return err
	}

	expandCache := &expandCache{}

	notify := func(sub *apps.Subscription) error {
		call := sub.Call
		if call == nil {
			return errors.New("nothing to call")
		}
		app, err := p.store.App.Get(sub.AppID)
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

func (p *proxy) upstreamForApp(app *apps.App) (upstream.Upstream, error) {
	m, err := p.store.Manifest.Get(app.AppID)
	if err != nil {
		return nil, err
	}

	switch m.Type {
	case apps.AppTypeHTTP:
		return uphttp.NewUpstream(m.RootURL, app.Secret), nil

	case apps.AppTypeAWS:
		return upawslambda.NewUpstream(app, m, p.aws), nil

	// case apps.AppTypeBuiltin:
	// 	if len(p.builtinProvisionedApps) == 0 {
	// 		return nil, errors.Errorf("builtin app not found: %s", app.Manifest.AppID)
	// 	}
	// 	up := p.builtinProvisionedApps[app.Manifest.AppID]
	// 	if up == nil {
	// 		return nil, errors.Errorf("builtin app not found: %s", app.Manifest.AppID)
	// 	}
	// 	return up, nil

	default:
		return nil, errors.Errorf("not a valid app type: %s", m.Type)
	}
}

func WriteCallError(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(apps.NewErrorCallResponse(err))
}
