// Copyright (c) 2020-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import "net/http"

type Proxy interface {
	GetBindings(*Context) ([]*Binding, error)
	Call(adminAccessToken string, _ *Call) (*Call, *CallResponse)
	Notify(*Context, Subject) error

	StartOAuthConnect(userID string, _ AppID, callOnComplete *Call) (connectURL string, _ error)
	HandleOAuth(http.ResponseWriter, *http.Request)

	ProvisionBuiltIn(AppID, Upstream)
}
