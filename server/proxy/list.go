// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package proxy

import (
	"strings"

	"github.com/mattermost/mattermost-plugin-apps/apps"
)

func (p *proxy) ListInstalledApps() map[apps.AppID]*apps.App {
	return p.store.App.List()
}

func (p *proxy) ListMarketplaceApps(filter string) map[apps.AppID]*apps.MarketplaceApp {
	out := map[apps.AppID]*apps.MarketplaceApp{}

	for appID, m := range p.store.Manifest.List() {
		if !appMatchesFilter(m, filter) {
			continue
		}
		marketApp := &apps.MarketplaceApp{
			Manifest: m,
		}
		app, _ := p.store.App.Get(appID)
		if app != nil {
			marketApp.Installed = true
			marketApp.Enabled = !app.Disabled
		}

		out[appID] = marketApp
	}

	return out
}

// Copied from Mattermost Server
func appMatchesFilter(manifest *apps.Manifest, filter string) bool {
	filter = strings.TrimSpace(strings.ToLower(filter))

	if filter == "" {
		return true
	}

	if strings.ToLower(string(manifest.AppID)) == filter {
		return true
	}

	if strings.Contains(strings.ToLower(manifest.DisplayName), filter) {
		return true
	}

	if strings.Contains(strings.ToLower(manifest.Description), filter) {
		return true
	}

	return false
}
