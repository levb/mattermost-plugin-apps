// Copyright (c) 2020-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package store

import (
	pluginapi "github.com/mattermost/mattermost-plugin-api"
	"github.com/mattermost/mattermost-plugin-apps/server/api"
)

const prefixSubs = "sub_"

type Store struct {
	builtinInstalledApps map[api.AppID]*api.App

	mm   *pluginapi.Client
	conf api.Configurator
}

var _ api.Store = (*Store)(nil)

func NewStore(mm *pluginapi.Client, conf api.Configurator) *Store {
	return &Store{
		builtinInstalledApps: map[api.AppID]*api.App{},
		mm:                   mm,
		conf:                 conf,
	}
}
