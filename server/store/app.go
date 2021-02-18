// Copyright (c) 2020-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package store

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/server/config"
	"github.com/mattermost/mattermost-plugin-apps/server/utils"
	"github.com/pkg/errors"
)

type App interface {
	config.Configurable

	Get(apps.AppID) (*apps.App, error)
	List() map[apps.AppID]*apps.App
	Store(*apps.App) error
	Delete(apps.AppID) error
	AddBuiltin(*apps.App)
}

type appStore struct {
	*Service

	installed map[apps.AppID]*apps.App
	builtin   map[apps.AppID]*apps.App
}

var _ App = (*appStore)(nil)

func (s *appStore) Configure(conf config.Config) error {
	s.installed = map[apps.AppID]*apps.App{}

	for id, key := range conf.InstalledApps {
		var app *apps.App
		err := s.mm.KV.Get(prefixInstalledApp+key, &app)
		if err != nil {
			s.mm.Log.Error(
				fmt.Sprintf("failed to load app %s: %s", id, err.Error()))
		}
		if app == nil {
			s.mm.Log.Error(
				fmt.Sprintf("failed to load app %s: key %s not found", id, prefixInstalledApp+key))
		}

		s.installed[apps.AppID(id)] = app
	}

	return nil
}

func (s *appStore) Get(appID apps.AppID) (*apps.App, error) {
	app, ok := s.installed[appID]
	if ok {
		return app, nil
	}
	app, ok = s.builtin[appID]
	if ok {
		return app, nil
	}
	return nil, utils.ErrNotFound
}

func (s *appStore) List() map[apps.AppID]*apps.App {
	out := map[apps.AppID]*apps.App{}
	for appID, app := range s.installed {
		out[appID] = app
	}

	for appID, app := range s.builtin {
		_, ok := s.installed[appID]
		if !ok {
			out[appID] = app
		}
	}
	return out
}

func (s *appStore) Store(app *apps.App) error {
	_, ok := s.installed[app.AppID]
	if ok {
		return errors.Errorf("failed to store: builtin app %s is read-only")
	}

	conf := s.conf.Get()
	prevSHA := conf.InstalledApps[string(app.AppID)]

	data, err := json.Marshal(app)
	if err != nil {
		return err
	}
	sha := fmt.Sprintf("%x", sha1.Sum(data))
	if sha == prevSHA {
		return nil
	}

	_, err = s.mm.KV.Set(prefixInstalledApp+sha, app)
	if err != nil {
		return err
	}
	updated := map[string]string{}
	for k, v := range conf.LocalManifests {
		updated[k] = v
	}
	updated[string(app.AppID)] = sha
	sc := *conf.StoredConfig
	sc.InstalledApps = updated
	err = s.conf.StoreConfig(&sc)
	if err != nil {
		return err
	}

	err = s.mm.KV.Delete(prefixInstalledApp + prevSHA)
	if err != nil {
		return err
	}

	return nil
}

func (s *appStore) Delete(appID apps.AppID) error {
	_, ok := s.installed[appID]
	if ok {
		return errors.Errorf("failed to delete: builtin app %s is read-only")
	}

	conf := s.conf.Get()
	sha, ok := conf.InstalledApps[string(appID)]
	if !ok {
		return utils.ErrNotFound
	}

	err := s.mm.KV.Delete(prefixInstalledApp + sha)
	if err != nil {
		return err
	}
	updated := map[string]string{}
	for k, v := range conf.InstalledApps {
		updated[k] = v
	}
	delete(updated, string(appID))
	sc := *conf.StoredConfig
	sc.InstalledApps = updated

	return s.conf.StoreConfig(&sc)
}

// AddBuiltinApp is not synchronized and should only be used at the plugin
// initialization time, to "register" builtin apps.
func (s *appStore) AddBuiltin(app *apps.App) {
	s.builtin[app.AppID] = app
}
