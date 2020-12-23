package proxy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mattermost/mattermost-plugin-api/experimental/bot/logger"
	"github.com/mattermost/mattermost-plugin-apps/server/utils/httputils"
	"github.com/pkg/errors"

	"golang.org/x/oauth2"

	"github.com/mattermost/mattermost-plugin-apps/server/api"
	"github.com/mattermost/mattermost-plugin-apps/server/utils/oauther"
)

func (p *Proxy) newMattermostOAuthenticator(app *api.App) oauther.OAuther {
	return oauther.NewFromClient(p.mm,
		*p.newMattermostOAuthConfig(app),
		p.finishMattermostOAuthConnect,
		logger.NewNilLogger(), // TODO replace with a real logger
		oauther.OAuthURL(oauther.DefaultOAuthURL+"/"+string(app.Manifest.AppID)),
		oauther.StorePrefix("mm_oauth_"))
}

func (p *Proxy) newMattermostOAuthConfig(app *api.App) *oauth2.Config {
	config := p.conf.GetConfig()
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

func (p *Proxy) HandleOAuth(w http.ResponseWriter, req *http.Request) {
	// the URL Path is .../oauth2/AppID/[connect|complete|...]
	fmt.Printf("<><> proxy.HandleOAuth 1: %s\n", req.URL)
	ss := strings.Split(req.URL.Path, "/")
	if len(ss) < 3 {
		httputils.WriteBadRequestError(w, errors.New("invalid path, can not extract AppID"))
		return
	}
	appID := api.AppID(ss[len(ss)-2])

	app, err := p.store.LoadApp(appID)
	if err != nil {
		httputils.WriteInternalServerError(w, err)
		return
	}

	p.newMattermostOAuthenticator(app).ServeHTTP(w, req)
}

func (p *Proxy) StartOAuthConnect(userID string, appID api.AppID, callOnComplete *api.Call) (string, error) {
	fmt.Printf("<><> startMattermostOAuthConnect 1: %v\n", userID)

	app, err := p.store.LoadApp(appID)
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

func (p *Proxy) finishMattermostOAuthConnect(userID string, token oauth2.Token, payload []byte) {
	call, err := api.UnmarshalCallFromData(payload)
	if err != nil {
		return
	}

	// TODO: figure out what to do with the CallResponse
	_, _ = p.Call("", call)
}
