// Copyright (c) 2020-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"net/http"

	"github.com/mattermost/mattermost-plugin-apps/apps"
)

type Proxy interface {
	GetBindings(*apps.Context) ([]*apps.Binding, error)
	//TODO: <><> get rid of sessionToken, should be in Context?
	Call(sessionToken string, call *apps.Call) (*apps.Call, *apps.CallResponse)
	Notify(cc *apps.Context, subj apps.Subject) error

	StartOAuthConnect(userID string, _ apps.AppID, callOnComplete *apps.Call) (connectURL string, _ error)
	HandleOAuth(http.ResponseWriter, *http.Request)

	ProvisionBuiltIn(apps.AppID, Upstream)
}
