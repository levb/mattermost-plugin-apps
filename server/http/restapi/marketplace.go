package restapi

import (
	"net/http"

	"github.com/mattermost/mattermost-plugin-apps/server/utils/httputils"
)

func (a *restapi) handleGetMarketplace(w http.ResponseWriter, req *http.Request, actingUserID string) {
	filter := req.URL.Query().Get("filter")
	marketplaceApps := a.proxy.ListMarketplaceApps(filter)
	httputils.WriteJSON(w, marketplaceApps)
}
