// Copyright (c) 2020-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package store

import (
	"crypto/sha1" // nolint:gosec
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/server/config"
	"github.com/mattermost/mattermost-plugin-apps/server/utils"
	"github.com/mattermost/mattermost-plugin-apps/server/utils/httputils"
)

// The list of all locally registered manifests is stored in the config as a map[AppID]=>sha256(manifest).
// The Save method updates the config and triggers a refresh across all
type Manifest interface {
	config.Configurable

	Init() error
	List() map[apps.AppID]*apps.Manifest
	Get(apps.AppID) (*apps.Manifest, error)
	StoreLocal(*apps.Manifest) error
	DeleteLocal(apps.AppID) error
}

type manifestStore struct {
	*Service

	global map[apps.AppID]*apps.Manifest
	local  map[apps.AppID]*apps.Manifest
}

var _ Manifest = (*manifestStore)(nil)

func (s *manifestStore) Init() error {
	bundlePath, err := s.mm.System.GetBundlePath()
	if err != nil {
		return errors.Wrap(err, "can't get bundle path")
	}
	assetPath := filepath.Join(bundlePath, "assets")
	f, err := os.Open(filepath.Join(assetPath, config.ManifestsFile))
	if err != nil {
		return errors.Wrap(err, "failed to load global list of available apps")
	}
	defer f.Close()

	return s.init(f, assetPath)
}

func (s *manifestStore) init(manifestsFile io.Reader, assetPath string) error {
	s.global = map[apps.AppID]*apps.Manifest{}

	// Read in the marketplace-listed manifests from S3, as per versions
	// indicated in apps.json. apps.json file contains a map of AppID->manifest
	// S3 filename (the bucket comes from the config)
	manifestLocations := map[apps.AppID]string{}
	err := json.NewDecoder(manifestsFile).Decode(&manifestLocations)
	if err != nil {
		return err
	}

	var data []byte
	conf := s.conf.Get()
	for appID, loc := range manifestLocations {
		parts := strings.SplitN(loc, ":", 2)
		switch {
		case len(parts) == 1:
			data, err = s.getFromS3(conf.AWSManifestBucket, appID, apps.AppVersion(parts[0]))
		case len(parts) == 2 && parts[0] == "s3":
			data, err = s.getFromS3(conf.AWSManifestBucket, appID, apps.AppVersion(parts[1]))
		case len(parts) == 2 && parts[0] == "file":
			data, err = ioutil.ReadFile(filepath.Join(assetPath, parts[1]))
		case len(parts) == 2 && (parts[0] == "http" || parts[0] == "https"):
			data, err = httputils.GetFromURL(loc)
		default:
			return errors.Errorf("failed to load global manifest for %s: %s is invalid", string(appID), loc)
		}
		if err != nil {
			return errors.Wrapf(err, "failed to load global manifest for %s", string(appID))
		}

		var m apps.Manifest
		err = json.Unmarshal(data, &m)
		if err != nil {
			return err
		}
		if m.AppID != appID {
			return errors.Errorf("missmatched app ids while getting manifest %s != %s", m.AppID, appID)
		}
		err = validateManifest(&m)
		if err != nil {
			return err
		}
		s.global[appID] = &m
	}

	return nil
}

func (s *manifestStore) Configure(conf config.Config) error {
	s.local = map[apps.AppID]*apps.Manifest{}

	for id, key := range conf.LocalManifests {
		var m *apps.Manifest
		err := s.mm.KV.Get(prefixLocalManifest+key, &m)
		if err != nil {
			s.mm.Log.Error(
				fmt.Sprintf("failed to load local manifest for %s: %s", id, err.Error()))
		}
		if m == nil {
			s.mm.Log.Error(
				fmt.Sprintf("failed to load local manifest for %s: not found", id))
		}

		s.local[apps.AppID(id)] = m
	}

	return nil
}

func (s *manifestStore) Get(appID apps.AppID) (*apps.Manifest, error) {
	m, ok := s.global[appID]
	if ok {
		return m, nil
	}
	m, ok = s.local[appID]
	if ok {
		return m, nil
	}
	return nil, utils.ErrNotFound
}

func (s *manifestStore) List() map[apps.AppID]*apps.Manifest {
	out := map[apps.AppID]*apps.Manifest{}
	for id, m := range s.global {
		out[id] = m
	}

	for id, m := range s.local {
		_, ok := s.global[id]
		if !ok {
			out[id] = m
		}
	}
	return out
}

func (s *manifestStore) StoreLocal(m *apps.Manifest) error {
	conf := s.conf.Get()
	prevSHA := conf.LocalManifests[string(m.AppID)]

	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	sha := fmt.Sprintf("%x", sha1.Sum(data)) // nolint:gosec
	if sha == prevSHA {
		return nil
	}

	_, err = s.mm.KV.Set(prefixLocalManifest+sha, m)
	if err != nil {
		return err
	}
	updated := map[string]string{}
	for k, v := range conf.LocalManifests {
		updated[k] = v
	}
	updated[string(m.AppID)] = sha
	sc := *conf.StoredConfig
	sc.LocalManifests = updated
	err = s.conf.StoreConfig(&sc)
	if err != nil {
		return err
	}

	err = s.mm.KV.Delete(prefixLocalManifest + prevSHA)
	if err != nil {
		return err
	}

	return nil
}

func (s *manifestStore) DeleteLocal(appID apps.AppID) error {
	conf := s.conf.Get()
	sha := conf.LocalManifests[string(appID)]

	err := s.mm.KV.Delete(prefixLocalManifest + sha)
	if err != nil {
		return err
	}
	updated := map[string]string{}
	for k, v := range conf.LocalManifests {
		updated[k] = v
	}
	delete(updated, string(appID))
	sc := *conf.StoredConfig
	sc.LocalManifests = updated

	return s.conf.StoreConfig(&sc)
}

// GetManifest returns a manifest file for an app from the S3
func (s *manifestStore) getFromS3(bucket string, appID apps.AppID, version apps.AppVersion) ([]byte, error) {
	name := fmt.Sprintf("manifest_%s_%s", appID, version)
	data, err := s.aws.Client().GetS3(bucket, name)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to download manifest %s", name)
	}
	return data, nil
}

func validateManifest(m *apps.Manifest) error {
	if m.AppID == "" {
		return errors.New("empty AppID")
	}
	if !m.Type.IsValid() {
		return errors.Errorf("invalid type: %s", m.Type)
	}

	if m.Type == apps.AppTypeHTTP {
		_, err := url.Parse(m.RootURL)
		if err != nil {
			return errors.Wrapf(err, "invalid manifest URL %q", m.RootURL)
		}
	}
	return nil
}
