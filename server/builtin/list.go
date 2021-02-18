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
	all := a.proxy.ListInstalledApps()

	// txt := md.MD(``)
	txt := md.MD("| Name  | Type | OAuth2 | Bot | Locations | Permissions |\n")
	txt += md.MD("| :-- |:-----| :----- | :-- | :-------- | :---------- |\n")
	for _, app := range all {
		txt += md.Markdownf("|[%s](%s) (%s)|%s|%s|%s|%s|%s|\n",
			app.DisplayName, app.HomepageURL, app.AppID, app.Type, app.OAuth2ClientID, app.BotUserID, app.GrantedLocations, app.GrantedPermissions)
	}

	return apps.NewCallResponse(txt, nil, nil)
}
