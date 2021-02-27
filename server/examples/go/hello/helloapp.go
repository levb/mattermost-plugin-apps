package hello

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
	PathSubmitSurvey             = "/survey-submit"
)

type HelloApp struct{}

func NewHelloApp() *HelloApp {
	return &HelloApp{}
}
