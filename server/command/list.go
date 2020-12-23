// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/mattermost/mattermost-plugin-apps/server/utils/md"
	"github.com/mattermost/mattermost-server/v5/model"
)

func (s *service) executeList(params *params) (*model.CommandResponse, error) {
	apps, err := s.api.Admin.ListApps()
	if err != nil {
		return errorOut(params, err)
	}

	txt := md.MD(`
	| ID  | Type | OAuth2 | Bot | Locations | Permissions |
	| :-- |:-----| :----- | :-- | :-------- | :---------- |
	`)

	for _, app := range apps {
		txt += md.Markdownf(`|%s|%s|%s|%s|%s|%s|`,
			app.Manifest.AppID, app.Manifest.Type, app.OAuth2ClientID, app.BotUserID, app.GrantedLocations, app.GrantedPermissions)
	}

	return out(params, txt)
}
