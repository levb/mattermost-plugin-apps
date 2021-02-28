package proxy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"

	"github.com/mattermost/mattermost-plugin-api/experimental/bot/logger"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/server/utils/httputils"
	"github.com/mattermost/mattermost-plugin-apps/server/utils/oauther"
)

func (p *proxy) newMattermostOAuthenticator(app *apps.App) oauther.OAuther {
	return oauther.NewFromClient(p.mm,
		*p.newMattermostOAuthConfig(app),
		p.finishMattermostOAuthConnect,
		logger.NewNilLogger(), // TODO replace with a real logger
		oauther.OAuthURL(oauther.DefaultOAuthURL+"/"+string(app.AppID)),
		oauther.StorePrefix("mm_oauth_"))
}

func (p *proxy) newMattermostOAuthConfig(app *apps.App) *oauth2.Config {
	config := p.conf.Get()
	return &oauth2.Config{
		ClientID:     app.OAuth2ClientID,
		ClientSecret: app.OAuth2ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.MattermostSiteURL + "/oauth/authorize",
			TokenURL: config.MattermostSiteURL + "/oauth/access_token",
		},
		// RedirectURL: - not needed, OAuther will configure
		// TODO: Scopes
	}
}

func (p *proxy) HandleOAuth(w http.ResponseWriter, req *http.Request) {
	// the URL Path is .../oauth2/AppID/[connect|complete|...]
	fmt.Printf("<><> proxy.HandleOAuth 1: %s\n", req.URL)
	ss := strings.Split(req.URL.Path, "/")
	if len(ss) < 3 {
		httputils.WriteBadRequestError(w, errors.New("invalid path, can not extract AppID"))
		return
	}
	appID := apps.AppID(ss[len(ss)-2])

	app, err := p.store.App.Get(appID)
	if err != nil {
		httputils.WriteInternalServerError(w, err)
		return
	}

	p.newMattermostOAuthenticator(app).ServeHTTP(w, req)
}

func (p *proxy) StartOAuthConnect(userID string, appID apps.AppID, callOnComplete *apps.Call) (string, error) {
	fmt.Printf("<><> startMattermostOAuthConnect 1: %v\n", userID)

	app, err := p.store.App.Get(appID)
	if err != nil {
		return "", err
	}
	oauth := p.newMattermostOAuthenticator(app)

	err = oauth.Deauthorize(userID)
	if err != nil {
		return "", err
	}
	state, err := json.Marshal(callOnComplete)
	if err != nil {
		return "", err
	}
	err = oauth.AddPayload(userID, state)
	if err != nil {
		return "", err
	}
	return oauth.GetConnectURL(), nil
}

func (p *proxy) finishMattermostOAuthConnect(userID string, token oauth2.Token, payload []byte) {
	call, err := apps.UnmarshalCallFromData(payload)
	if err != nil {
		return
	}

	// TODO: figure out what to do with the CallResponse
	_, _ = p.Call("", call)
}
