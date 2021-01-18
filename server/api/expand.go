// Copyright (c) 2020-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

import (
	"strings"
)

type ExpandLevel string

// Expand specifies the data dependencies of the Call that should be provided in
// its Context at the invocation time.
//
// Example:
// ```json
//	"call":{
//		"url": "/some-path",
//		"expand":{
//			"acting_user_access_token":"required",
//			"channel":"summary",
//			"post":"summary",
//		}
//	}
// ```
type Expand struct {
	App        ExpandLevel `json:"app,omitempty"`
	ActingUser ExpandLevel `json:"acting_user,omitempty"`

	// ActingUserAccessToken instruct the proxy to include OAuth2 access token
	// in the request. If the token is not available or is invalid, the user is
	// directed to the OAuth2 flow, and the Call is executed upon completion.
	ActingUserAccessToken ExpandLevel `json:"acting_user_access_token,omitempty"`

	// AdminAccessToken instructs the proxy to include an admin access token.
	AdminAccessToken ExpandLevel `json:"admin_access_token,omitempty"`

	Channel   ExpandLevel `json:"channel,omitempty"`
	Mentioned ExpandLevel `json:"mentioned,omitempty"`
	Post      ExpandLevel `json:"post,omitempty"`
	RootPost  ExpandLevel `json:"root_post,omitempty"`
	Team      ExpandLevel `json:"team,omitempty"`
	User      ExpandLevel `json:"user,omitempty"`
}

const (
	ExpandDefault  = ExpandLevel("")
	ExpandNone     = ExpandLevel("none")
	ExpandAll      = ExpandLevel("all")
	ExpandSummary  = ExpandLevel("digest")
	ExpandRequired = ExpandLevel("required")
	ExpandOptional = ExpandLevel("optional")
)

func (el ExpandLevel) Any() bool {
	return el != "" && el != ExpandNone
}

func (el ExpandLevel) IsRequired() bool {
	return el.Contains(ExpandRequired)
}

func (el ExpandLevel) Contains(level ExpandLevel) bool {
	ss := strings.Split(string(el), ",")
	for _, s := range ss {
		if strings.TrimSpace(s) == string(level) {
			return true
		}
	}
	return false
}
