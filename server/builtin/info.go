// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package builtin

import (
	"github.com/mattermost/mattermost-plugin-apps/server/api"
	"github.com/mattermost/mattermost-plugin-apps/server/utils/md"
)

func (a *App) info(call *api.Call) (md.MD, error) {
	conf := a.API.Configurator.GetConfig()

	resp := md.Markdownf("Mattermost Cloud Apps plugin version: %s, "+
		"[%s](https://github.com/mattermost/%s/commit/%s), built %s\n",
		conf.Version,
		conf.BuildHashShort,
		api.Repository,
		conf.BuildHash,
		conf.BuildDate)
	return resp, nil
}
