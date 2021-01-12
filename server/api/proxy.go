// Copyright (c) 2020-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import "net/http"

type Proxy interface {
	GetBindings(*Context) ([]*Binding, error)
	Call(SessionToken, *Call) (*Call, *CallResponse)
	Notify(cc *Context, subj Subject) error

	StartOAuthConnect(userID string, appID AppID, callOnComplete *Call) (connectURL string, err error)
	HandleOAuth(http.ResponseWriter, *http.Request)

	ProvisionBuiltIn(AppID, Upstream)
}
