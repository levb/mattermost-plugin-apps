package hello

import (
	"github.com/mattermost/mattermost-plugin-apps/server/api"
)

const (
	fieldUserID   = "userID"
	fieldMessage  = "message"
	fieldResponse = "response"
)

const (
	PathPostAsUser               = "/post-as-user"
	PathSendSurvey               = "/send"
	PathSendSurveyCommandToModal = "/send-command-modal"
	PathSendSurveyModal          = "/send-modal"
	PathSubscribeChannel         = "/subscribe"
	PathSurvey                   = "/survey"
	PathUserJoinedChannel        = "/user-joined-channel"
)

type HelloApp struct {
	API *api.Service
}

func NewHelloApp(api *api.Service) *HelloApp {
	return &HelloApp{api}
}
