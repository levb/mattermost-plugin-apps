// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package admin

import (
	"github.com/mattermost/mattermost-plugin-apps/server/api"
)

func (adm *Admin) ListApps() ([]*api.App, error) {
	return adm.store.ListApps(), nil
}
