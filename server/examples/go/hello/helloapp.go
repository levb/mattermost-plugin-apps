package hello

import "github.com/mattermost/mattermost-plugin-apps/server/config"

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
	Conf config.Service
}

func NewHelloApp(conf config.Service) *HelloApp {
	return &HelloApp{
		Conf: conf,
	}
}
