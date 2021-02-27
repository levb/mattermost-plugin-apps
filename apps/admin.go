// Copyright (c) 2020-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package apps

type InInstallApp struct {
	Manifest         *Manifest `json:"manifest"`
	OAuth2TrustedApp bool      `json:"oauth2_trusted_app,omitempty"`
	Secret           string    `json:"secret,omitempty"`
	Force            bool      `json:"force,omitempty"`
}
