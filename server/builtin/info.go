// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package builtin

import (
	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/server/config"
	"github.com/mattermost/mattermost-plugin-apps/server/utils/md"
)

func (a *builtinApp) infoForm(c *apps.Call) *apps.CallResponse {
	return &apps.CallResponse{
		Type: apps.CallResponseTypeForm,
		Form: &apps.Form{
			Title: "Apps proxy info",
			Call: &apps.Call{
				URL: PathInfo,
			},
		},
	}
}

func (a *builtinApp) info(call *apps.Call) *apps.CallResponse {
	conf := a.conf.Get()
	resp := md.Markdownf("Mattermost Cloud Apps plugin version: %s, "+
		"[%s](https://github.com/mattermost/%s/commit/%s), built %s\n",
		conf.Version,
		conf.BuildHashShort,
		config.Repository,
		conf.BuildHash,
		conf.BuildDate)

	return apps.NewCallResponse(resp, nil, nil)
}
