package hello

import (
	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/apps/mmclient"
)

func (h *HelloApp) UserJoinedChannel(call *apps.Call) *apps.CallResponse {
	go func() {
		bot := mmclient.AsBot(call.Context)

		_ = sendSurvey(bot, call.Context.UserID, "welcome to channel")
	}()
	return apps.NewCallResponse("ok", nil, nil)
}
