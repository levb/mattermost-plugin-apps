// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package proxy

import (
	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/server/utils/md"
)

func (p *proxy) InstallApp(cc *apps.Context, in *apps.InInstallApp) (*apps.App, md.MD, error) {
	// 	if !p.mm.User.HasPermissionTo(cc.ActingUserID, model.PERMISSION_MANAGE_SYSTEM) {
	// 		return nil, "", errors.New("forbidden")
	// 	}

	// 	app := in.App
	// 	if app.AppID == "" {
	// 		return nil, "", errors.New("app ID must not be empty")
	// 	}
	// 	_, err := adm.store.App().Get(app.Manifest.AppID)
	// 	switch {
	// 	case err == utils.ErrNotFound:

	// 	case err == nil && !in.Force:
	// 		return nil, "", errors.Errorf("app %s already provisioned, use Force to overwrite", app.Manifest.AppID)

	// 	default:
	// 		return nil, "", err
	// 	}

	// 	bot, token, err := adm.ensureBot(app.Manifest, cc.ActingUserID, cc.AdminAccessToken)
	// 	if err != nil {
	// 		return nil, "", err
	// 	}
	// 	app.BotUserID = bot.UserId
	// 	app.BotUsername = bot.Username
	// 	app.BotAccessToken = token.Token

	// 	conf := adm.conf.GetConfig()
	// 	client := model.NewAPIv4Client(conf.MattermostSiteURL)
	// 	client.SetToken(cc.AdminAccessToken)

	// 	if app.GrantedPermissions.Contains(apps.PermissionActAsUser) {
	// 		var oAuthApp *model.OAuthApp
	// 		oAuthApp, err = adm.ensureOAuthApp(app.Manifest, app.OAuth2TrustedApp, cc.ActingUserID, string(cc.AdminAccessToken))
	// 		if err != nil {
	// 			return nil, "", err
	// 		}
	// 		app.OAuth2ClientID = oAuthApp.Id
	// 		app.OAuth2ClientSecret = oAuthApp.ClientSecret

	// 		// Installed app is automatically enabled, since config is done in the installation process
	// 		if app.Status == "" || app.Status == apps.AppStatusRegistered {
	// 			app.Status = apps.AppStatusInstalled
	// 		}
	// 	}

	// 	err = adm.store.App().Save(app)
	// 	if err != nil {
	// 		return nil, "", err
	// 	}

	// 	install := app.Manifest.OnInstall
	// 	if install == nil {
	// 		install = apps.DefaultInstallCall
	// 	}
	// 	install.Context = cc
	// 	install.Context.ExpandedContext = apps.ExpandedContext{}

	// 	_, cr := adm.proxy.Call(cc.AdminAccessToken, install)
	// 	if cr.Type == apps.CallResponseTypeError {
	// 		return nil, "", errors.Wrap(cr, "install failed")
	// 	}

	// 	return app, cr.Markdown, nil
	// }

	// func (adm *Admin) ensureBot(manifest *apps.Manifest, actingUserID, sessionToken string) (*model.Bot, *model.UserAccessToken, error) {
	// 	conf := adm.conf.GetConfig()
	// 	client := model.NewAPIv4Client(conf.MattermostSiteURL)
	// 	client.SetToken(sessionToken)

	// 	bot := &model.Bot{
	// 		Username:    strings.ToLower(string(manifest.AppID)),
	// 		DisplayName: manifest.DisplayName,
	// 		Description: fmt.Sprintf("Bot account for `%s` App.", manifest.DisplayName),
	// 	}

	// 	var fullBot *model.Bot
	// 	user, _ := client.GetUserByUsername(bot.Username, "")
	// 	if user == nil {
	// 		var response *model.Response
	// 		fullBot, response = client.CreateBot(bot)

	// 		if response.StatusCode != http.StatusCreated {
	// 			if response.Error != nil {
	// 				return nil, nil, response.Error
	// 			}
	// 			return nil, nil, errors.New("could not create bot")
	// 		}
	// 	} else {
	// 		if !user.IsBot {
	// 			return nil, nil, errors.New("a user already owns the bot username")
	// 		}

	// 		fullBot = model.BotFromUser(user)
	// 		if fullBot.DeleteAt != 0 {
	// 			var response *model.Response
	// 			fullBot, response = client.EnableBot(fullBot.UserId)
	// 			if response.StatusCode != http.StatusOK {
	// 				if response.Error != nil {
	// 					return nil, nil, response.Error
	// 				}
	// 				return nil, nil, errors.New("could not enable bot")
	// 			}
	// 		}
	// 	}

	// 	token, response := client.CreateUserAccessToken(fullBot.UserId, "Mattermost App Token")
	// 	if response.StatusCode != http.StatusOK {
	// 		if response.Error != nil {
	// 			return nil, nil, response.Error
	// 		}
	// 		return nil, nil, fmt.Errorf("could not create token, status code = %v", response.StatusCode)
	// 	}

	// 	_ = adm.mm.Post.DM(fullBot.UserId, actingUserID, &model.Post{
	// 		Message: fmt.Sprintf("Provisioned bot account @%s (`%s`).",
	// 			fullBot.Username, fullBot.UserId),
	// 	})

	// 	return fullBot, token, nil
	// }

	// func (adm *Admin) ensureOAuthApp(manifest *apps.Manifest, noUserConsent bool, actingUserID, sessionToken string) (*model.OAuthApp, error) {
	// 	app, err := adm.store.App().Get(manifest.AppID)
	// 	if err != nil && err != utils.ErrNotFound {
	// 		return nil, err
	// 	}

	// 	conf := adm.conf.GetConfig()
	// 	client := model.NewAPIv4Client(conf.MattermostSiteURL)
	// 	client.SetToken(sessionToken)

	// 	if app.OAuth2ClientID != "" {
	// 		oauthApp, response := client.GetOAuthApp(app.OAuth2ClientID)
	// 		if response.StatusCode == http.StatusOK && response.Error == nil {
	// 			_ = adm.mm.Post.DM(app.BotUserID, actingUserID, &model.Post{
	// 				Message: fmt.Sprintf("Using existing OAuth2 App `%s`.", oauthApp.Id),
	// 			})

	// 			return oauthApp, nil
	// 		}
	// 	}

	// 	oauth2CallbackURL := fmt.Sprintf("%s%s/%s%s",
	// 		adm.conf.GetConfig().PluginURL, api.OAuth2Path, string(manifest.AppID), oauther.CompletePath)

	// 	// For the POC this should work, but for the final product I would opt for a RPC method to register the App
	// 	oauthApp, response := client.CreateOAuthApp(&model.OAuthApp{
	// 		CreatorId:    actingUserID,
	// 		Name:         manifest.DisplayName,
	// 		Description:  manifest.Description,
	// 		CallbackUrls: []string{oauth2CallbackURL},
	// 		Homepage:     manifest.HomepageURL,
	// 		IsTrusted:    noUserConsent,
	// 	})
	// 	if response.StatusCode != http.StatusCreated {
	// 		if response.Error != nil {
	// 			return nil, errors.Wrap(response.Error, "failed to create OAuth2 App")
	// 		}
	// 		return nil, errors.Errorf("failed to create OAuth2 App: received status code %v", response.StatusCode)
	// 	}

	// 	_ = adm.mm.Post.DM(app.BotUserID, actingUserID, &model.Post{
	// 		Message: fmt.Sprintf("Created OAuth2 App (`%s`). Callback URL: %s", oauthApp.Id, oauth2CallbackURL),
	// 	})

	// 	return oauthApp, nil

	return nil, "", nil
}
