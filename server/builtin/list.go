// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package builtin

import (
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
	all, err := a.API.Admin.ListApps()
	if err != nil {
		return apps.NewCallResponse("", nil, err)
	}

	// txt := md.MD(``)
	txt := md.MD("| ID  | Type | OAuth2 | Bot | Locations | Permissions |\n")
	txt += md.MD("| :-- |:-----| :----- | :-- | :-------- | :---------- |\n")
	for _, app := range all {
		// txt += md.Markdownf(`%s - %s - %s - %s`,
		// 	app.Manifest.AppID, app.Manifest.Type, app.GrantedLocations, app.GrantedPermissions)
		txt += md.Markdownf("|%s|%s|%s|%s|%s|%s|\n",
			app.Manifest.AppID, app.Manifest.Type, app.OAuth2ClientID, app.BotUserID, app.GrantedLocations, app.GrantedPermissions)
	}

	return apps.NewCallResponse(txt, nil, nil)
}
