// Copyright (c) 2020-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"github.com/mattermost/mattermost-plugin-apps/server/utils/md"
)

type Admin interface {
	ListApps() ([]*App, error)
	InstallApp(*Context, *InInstallApp) (*App, md.MD, error)
}

type InInstallApp struct {
	App   *App `json:"app"`
	Force bool `json:"force,omitempty"`
}
