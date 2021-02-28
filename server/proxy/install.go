// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package proxy

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/server/config"
	"github.com/mattermost/mattermost-plugin-apps/server/utils"
	"github.com/mattermost/mattermost-plugin-apps/server/utils/md"
	"github.com/mattermost/mattermost-plugin-apps/server/utils/oauther"
)

func (p *proxy) InstallApp(cc *apps.Context, in *apps.InInstallApp) (*apps.App, md.MD, error) {
	if !p.mm.User.HasPermissionTo(cc.ActingUserID, model.PERMISSION_MANAGE_SYSTEM) {
		return nil, "", errors.New("forbidden")
	}

	if in.Manifest.AppID == "" {
		return nil, "", errors.New("app ID must not be empty")
	}
	app, err := p.store.App.Get(in.Manifest.AppID)
	if err != nil {
		if err == utils.ErrNotFound {
			app = &apps.App{}
		} else {
			return nil, "", err
		}
	}
	app.Common = in.Manifest.Common
	app.GrantedLocations = in.Manifest.RequestedLocations
	app.GrantedPermissions = in.Manifest.RequestedPermissions
	app.Secret = in.Secret

	if app.BotAccessToken == "" {
		bot, token, botErr := p.ensureBot(in.Manifest, cc.ActingUserID, cc.AdminAccessToken)
		if botErr != nil {
			return nil, "", botErr
		}
		app.BotUserID = bot.UserId
		app.BotUsername = bot.Username
		app.BotAccessToken = token.Token
	}

	if (in.Manifest.RequestedPermissions.Contains(apps.PermissionActAsUser)) &&
		app.OAuth2ClientID == "" {
		var oAuthApp *model.OAuthApp
		oAuthApp, err = p.ensureOAuthApp(in.Manifest, in.OAuth2TrustedApp, cc.ActingUserID, cc.AdminAccessToken)
		if err != nil {
			return nil, "", err
		}
		app.OAuth2TrustedApp = in.OAuth2TrustedApp
		app.OAuth2ClientID = oAuthApp.Id
		app.OAuth2ClientSecret = oAuthApp.ClientSecret
	}
	app.Disabled = false

	err = p.store.App.Store(app)
	if err != nil {
		return nil, "", err
	}

	install := in.Manifest.OnInstall
	if install == nil {
		install = apps.DefaultInstallCall
	}
	cloneCC := *cc
	cloneCC.AppID = app.AppID
	cloneCC.ExpandedContext = apps.ExpandedContext{}
	install.Context = &cloneCC

	_, cr := p.Call(cc.AdminAccessToken, install)
	if cr.Type == apps.CallResponseTypeError {
		return nil, "", errors.Wrap(cr, "install failed")
	}

	return app, cr.Markdown, nil
}

func (p *proxy) ensureBot(manifest *apps.Manifest, actingUserID, sessionToken string) (*model.Bot, *model.UserAccessToken, error) {
	conf := p.conf.Get()
	client := model.NewAPIv4Client(conf.MattermostSiteURL)
	client.SetToken(sessionToken)

	bot := &model.Bot{
		Username:    strings.ToLower(string(manifest.AppID)),
		DisplayName: manifest.DisplayName,
		Description: fmt.Sprintf("Bot account for `%s` App.", manifest.DisplayName),
	}

	var fullBot *model.Bot
	user, _ := client.GetUserByUsername(bot.Username, "")
	if user == nil {
		var response *model.Response
		fullBot, response = client.CreateBot(bot)

		if response.StatusCode != http.StatusCreated {
			if response.Error != nil {
				return nil, nil, response.Error
			}
			return nil, nil, errors.New("could not create bot")
		}
	} else {
		if !user.IsBot {
			return nil, nil, errors.New("a user already owns the bot username")
		}

		fullBot = model.BotFromUser(user)
		if fullBot.DeleteAt != 0 {
			var response *model.Response
			fullBot, response = client.EnableBot(fullBot.UserId)
			if response.StatusCode != http.StatusOK {
				if response.Error != nil {
					return nil, nil, response.Error
				}
				return nil, nil, errors.New("could not enable bot")
			}
		}
	}

	token, response := client.CreateUserAccessToken(fullBot.UserId, "Mattermost App Token")
	if response.StatusCode != http.StatusOK {
		if response.Error != nil {
			return nil, nil, response.Error
		}
		return nil, nil, fmt.Errorf("could not create token, status code = %v", response.StatusCode)
	}

	_ = p.mm.Post.DM(fullBot.UserId, actingUserID, &model.Post{
		Message: fmt.Sprintf("Provisioned bot account @%s (`%s`).",
			fullBot.Username, fullBot.UserId),
	})

	return fullBot, token, nil
}

func (p *proxy) ensureOAuthApp(m *apps.Manifest, noUserConsent bool, actingUserID, sessionToken string) (*model.OAuthApp, error) {
	conf := p.conf.Get()
	client := model.NewAPIv4Client(conf.MattermostSiteURL)
	client.SetToken(sessionToken)

	oauth2CallbackURL := fmt.Sprintf("%s%s/%s%s",
		conf.PluginURL, config.OAuth2Path, m.AppID, oauther.CompletePath)

	// For the POC this should work, but for the final product I would opt for a RPC method to register the App
	oauthApp, response := client.CreateOAuthApp(&model.OAuthApp{
		CreatorId:    actingUserID,
		Name:         m.DisplayName,
		Description:  m.Description,
		CallbackUrls: []string{oauth2CallbackURL},
		Homepage:     m.HomepageURL,
		IsTrusted:    noUserConsent,
	})
	if response.StatusCode != http.StatusCreated {
		if response.Error != nil {
			return nil, errors.Wrap(response.Error, "failed to create OAuth2 App")
		}
		return nil, errors.Errorf("failed to create OAuth2 App: received status code %v", response.StatusCode)
	}

	return oauthApp, nil
}
