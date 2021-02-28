// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package proxy

import (
	"github.com/mattermost/mattermost-plugin-apps/apps"
)

func (p *proxy) UninstallApp(cc *apps.Context, appID apps.AppID) error {
	// // Call delete the function of the app
	// if err := adm.expandedCall(sessionToken, app, app.Manifest.OnUninstall, nil); err != nil {
	// 	return errors.Wrapf(err, "uninstall failed. appID - %s", app.Manifest.AppID)
	// }

	// // delete oauth app
	// conf := adm.conf.GetConfig()
	// client := model.NewAPIv4Client(conf.MattermostSiteURL)
	// client.SetToken(cc.AdminAccessToken)

	// if app.OAuth2ClientID != "" {
	// 	success, response := client.DeleteOAuthApp(app.OAuth2ClientID)
	// 	if !success || response.StatusCode != http.StatusNoContent {
	// 		return errors.Wrapf(response.Error, "failed to delete OAuth2 App - %s", app.Manifest.AppID)
	// 	}
	// }

	// // delete the bot account
	// if err := adm.mm.Bot.DeletePermanently(app.BotUserID); err != nil {
	// 	return errors.Wrapf(err, "can't delete bot account for App - %s", app.Manifest.AppID)
	// }

	// // delete app from proxy plugin, not removing the data
	// if err := adm.store.App().Delete(app); err != nil {
	// 	return errors.Wrapf(err, "can't delete app - %s", app.Manifest.AppID)
	// }

	// adm.mm.Log.Info("Uninstalled the app", "app_id", app.Manifest.AppID)

	return nil
}
