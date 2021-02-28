package restapi

import (
	"net/http"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/server/config"
	"github.com/mattermost/mattermost-plugin-apps/server/utils/httputils"
)

func (a *restapi) handleGetBindings(w http.ResponseWriter, req *http.Request, actingUserID string) {
	query := req.URL.Query()
	bindings, err := a.proxy.GetBindings(&apps.Context{
		TeamID:            query.Get(config.PropTeamID),
		ChannelID:         query.Get(config.PropChannelID),
		ActingUserID:      actingUserID,
		UserID:            actingUserID,
		PostID:            query.Get(config.PropPostID),
		MattermostSiteURL: a.conf.Get().MattermostSiteURL,
	})
	if err != nil {
		httputils.WriteInternalServerError(w, err)
		return
	}

	httputils.WriteJSON(w, bindings)
}
