// Copyright (c) 2020-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package proxy

import (
	"net/http"

	pluginapi "github.com/mattermost/mattermost-plugin-api"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/server/aws"
	"github.com/mattermost/mattermost-plugin-apps/server/config"
	"github.com/mattermost/mattermost-plugin-apps/server/store"
	"github.com/mattermost/mattermost-plugin-apps/server/upstream"
	"github.com/mattermost/mattermost-plugin-apps/server/utils/md"
)

type BuiltinApp interface {
	upstream.Upstream
	App() *apps.App
}

type Service interface {
	ListInstalledApps() map[apps.AppID]*apps.App
	ListMarketplaceApps(filter string) []*apps.MarketplaceApp

	InstallApp(*apps.Context, *apps.InInstallApp) (*apps.App, md.MD, error)
	UninstallApp(*apps.Context, apps.AppID) error
	// EnableApp(*apps.Context, apps.AppID) error
	// DisableApp(*apps.Context, apps.AppID) error

	GetBindings(*apps.Context) ([]*apps.Binding, error)
	//TODO: <><> get rid of sessionToken, should be in Context?
	Call(sessionToken string, call *apps.Call) (*apps.Call, *apps.CallResponse)
	Notify(cc *apps.Context, subj apps.Subject) error

	StartOAuthConnect(userID string, _ apps.AppID, callOnComplete *apps.Call) (connectURL string, _ error)
	HandleOAuth(http.ResponseWriter, *http.Request)

	AddBuiltin(BuiltinApp)
}

type proxy struct {
	// Manifests contains all relevant manifests. For V1, the entire list is
	// cached in memory, and loaded on startup.
	Manifests map[apps.AppID]*apps.Manifest

	// Built-in Apps are linked in Go and invoked directly. The list is
	// initialized on startup, and need not be synchronized. Built-in apps do
	// not need manifests.
	builtinUpstreams map[apps.AppID]upstream.Upstream

	mm    *pluginapi.Client
	conf  config.Service
	store *store.Service
	aws   aws.Service
}

var _ Service = (*proxy)(nil)

func NewService(mm *pluginapi.Client, aws aws.Service, conf config.Service, store *store.Service) *proxy {
	return &proxy{
		mm:    mm,
		conf:  conf,
		store: store,
		aws:   aws,
	}
}

func (p *proxy) AddBuiltin(builtinApp BuiltinApp) {
	app := builtinApp.App()

	p.store.App.AddBuiltin(app)

	if p.builtinUpstreams == nil {
		p.builtinUpstreams = map[apps.AppID]upstream.Upstream{}
	}
	p.builtinUpstreams[app.AppID] = builtinApp
}
