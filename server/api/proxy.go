// Copyright (c) 2020-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import "net/http"

type Proxy interface {
	GetBindings(*Context) ([]*Binding, error)
	Call(SessionToken, *Call) *CallResponse
	Notify(cc *Context, subj Subject) error
	HandleOAuth(w http.ResponseWriter, req *http.Request)

	ProvisionBuiltIn(AppID, Upstream)
}
