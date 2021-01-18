// Copyright (c) 2020-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package store

import (
	"github.com/mattermost/mattermost-plugin-apps/server/api"
	"github.com/mattermost/mattermost-plugin-apps/server/utils"
	"github.com/pkg/errors"
)

func (s *Store) ListApps() []*api.App {
	out := []*api.App{}
	for _, app := range s.builtinInstalledApps {
		out = append(out, app)
	}

	conf := s.conf.GetConfig()
	if conf.Apps == nil {
		return nil
	}
	for _, v := range conf.Apps {
		app := api.AppFromConfigMap(v)
		out = append(out, app)
	}
	return out
}

func (s *Store) LoadApp(appID api.AppID) (*api.App, error) {
	app := s.builtinInstalledApps[appID]
	if app != nil {
		return app, nil
	}

	conf := s.conf.GetConfig()
	if len(conf.Apps) == 0 {
		return nil, utils.ErrNotFound
	}
	v := conf.Apps[string(appID)]
	if v == nil {
		return nil, utils.ErrNotFound
	}
	return api.AppFromConfigMap(v), nil
}

func (s *Store) StoreApp(app *api.App) error {
	if s.builtinInstalledApps[app.Manifest.AppID] != nil {
		return errors.Errorf("failed to store app: %s is a builtin.", app.Manifest.AppID)
	}

	conf := s.conf.GetConfig()
	if len(conf.Apps) == 0 {
		conf.Apps = map[string]interface{}{}
	}

	conf.Apps[string(app.Manifest.AppID)] = app.ConfigMap()

	// Refresh the local config immediately, do not wait for the
	// OnConfigurationChange.
	err := s.conf.RefreshConfig(conf.StoredConfig)
	if err != nil {
		return err
	}

	return s.conf.StoreConfig(conf.StoredConfig)
}

// AddBuiltinApp is not synchronized and should only be used at the plugin
// initialization time, to "register" builtin apps.
func (s *Store) AddBuiltinApp(app *api.App) {
	s.builtinInstalledApps[app.Manifest.AppID] = app
}
