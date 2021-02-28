package restapi

import (
	"net/http"
	"sort"
	"strings"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/server/utils/httputils"
)

func (a *restapi) handleGetMarketplace(w http.ResponseWriter, req *http.Request, actingUserID string) {
	filter := req.URL.Query().Get("filter")
	out := []*apps.MarketplaceApp{}
	for _, mapp := range a.proxy.ListMarketplaceApps(filter) {
		out = append(out, mapp)
	}

	// Sort by display name, alphabetically.
	sort.SliceStable(out, func(i, j int) bool {
		return strings.ToLower(out[i].Manifest.DisplayName) <
			strings.ToLower(out[j].Manifest.DisplayName)
	})

	httputils.WriteJSON(w, out)
}
