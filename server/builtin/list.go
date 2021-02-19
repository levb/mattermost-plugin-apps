// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package builtin

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/server/utils/md"
)

func (a *App) listForm(c *apps.Call) (*apps.Form, error) {
	return &apps.Form{
		Title: "list Apps",
		Call: &apps.Call{
			URL: PathList,
		},
	}, nil
}

func (a *App) list(call *apps.Call) *apps.CallResponse {
	marketplaceApps := a.proxy.ListMarketplaceApps("")
	installedApps := a.proxy.ListInstalledApps()

	txt := md.MD("| Name | Status | Version | Account | Locations | Permissions |\n")
	txt += md.MD("| :-- |:-- | :-- | :-- | :-- | :-- |\n")

	for _, app := range installedApps {
		mapp := marketplaceApps[app.AppID]

		status := "Installed"
		if app.Disabled {
			status += ", Disabled"
		}
		if mapp == nil {
			status += ", Unlisted"
		}
		status += fmt.Sprintf(", type: `%s`", app.Type)

		version := string(app.Version)
		if mapp != nil && string(mapp.Manifest.Version) != version {
			version += fmt.Sprintf("(marketplace: %s)", mapp.Manifest.Version)
		}

		account := ""
		if app.BotUserID != "" {
			account += fmt.Sprintf("Bot: `%s`", app.BotUserID)
		}
		if app.OAuth2ClientID != "" {
			if account != "" {
				account += ", "
			}
			account += fmt.Sprintf("OAuth: `%s`", app.OAuth2ClientID)
		}
		name := fmt.Sprintf("[%s](%s) (%s)",
			app.DisplayName, app.HomepageURL, app.AppID)

		txt += md.Markdownf("|%s|%s|%s|%s|%s|%s|\n",
			name, status, version, account, app.GrantedLocations, app.GrantedPermissions)
	}

	for _, mapp := range marketplaceApps {
		_, ok := installedApps[mapp.Manifest.AppID]
		if ok {
			continue
		}
		version := string(mapp.Manifest.Version)
		name := fmt.Sprintf("[%s](%s) (%s)",
			mapp.Manifest.DisplayName, mapp.Manifest.HomepageURL, mapp.Manifest.AppID)
		txt += md.Markdownf("|%s|%s|%s|%s|%s|%s|\n",
			name, "", version, "", mapp.Manifest.RequestedLocations, mapp.Manifest.RequestedPermissions)
	}

	return apps.NewCallResponse(txt, nil, nil)
}
