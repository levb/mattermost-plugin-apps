// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package builtin

import (
	"github.com/mattermost/mattermost-plugin-apps/server/api"
	"github.com/mattermost/mattermost-plugin-apps/server/utils/md"
)

func (a *App) listForm(c *api.Call) (*api.Form, error) {
	return &api.Form{
		Title: "list Apps",
		Call: &api.Call{
			URL: PathList,
		},
	}, nil
}

func (a *App) list(call *api.Call) (md.MD, error) {
	apps, err := a.API.Admin.ListApps()
	if err != nil {
		return "", err
	}

	// txt := md.MD(``)
	txt := md.MD("| ID  | Type | OAuth2 | Bot | Locations | Permissions |\n")
	txt += md.MD("| :-- |:-----| :----- | :-- | :-------- | :---------- |\n")
	for _, app := range apps {
		// txt += md.Markdownf(`%s - %s - %s - %s`,
		// 	app.Manifest.AppID, app.Manifest.Type, app.GrantedLocations, app.GrantedPermissions)
		txt += md.Markdownf("|%s|%s|%s|%s|%s|%s|\n",
			app.Manifest.AppID, app.Manifest.Type, app.OAuth2ClientID, app.BotUserID, app.GrantedLocations, app.GrantedPermissions)
	}

	return txt, nil
}
