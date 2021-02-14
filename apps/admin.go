// Copyright (c) 2020-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package apps

type InInstallApp struct {
	App   *App `json:"app"`
	Force bool `json:"force,omitempty"`
}
