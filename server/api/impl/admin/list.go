// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package admin

import "github.com/mattermost/mattermost-plugin-apps/apps"

func (adm *Admin) ListApps() ([]*apps.App, error) {
	return adm.store.App().GetAll(), nil
}
