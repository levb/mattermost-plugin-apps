package hello

import (
	"github.com/mattermost/mattermost-plugin-apps/server/api"
	"github.com/pkg/errors"
)

func (h *HelloApp) Survey(call *api.Call) *api.CallResponse {
	switch call.Type {
	case api.CallTypeForm:
		return newSurveyFormResponse(call)

	case api.CallTypeSubmit:
		err := h.processSurvey(call)
		return api.NewCallResponse("ok", nil, err)

	default:
		return api.NewErrorCallResponse(errors.New("not supported"))
	}
}

func NewSurveyForm(message string) *api.Form {
	return &api.Form{
		Title:         "Emotional response survey",
		Header:        message,
		Footer:        "Let the world know!",
		SubmitButtons: fieldResponse,
		Fields: []*api.Field{
			{
				Name: fieldResponse,
				Type: api.FieldTypeStaticSelect,
				SelectStaticOptions: []api.SelectOption{
					{Label: "Like", Value: "like"},
					{Label: "Dislike", Value: "dislike"},
				},
			},
		},
	}
}

func newSurveyFormResponse(c *api.Call) *api.CallResponse {
	message := c.GetStringValue(fieldMessage, "default hello message")
	return &api.CallResponse{
		Type: api.CallResponseTypeForm,
		Form: NewSurveyForm(message),
	}
}

func (h *HelloApp) processSurvey(c *api.Call) error {
	// TODO post something; for embedded form - what do we do?
	return nil
}
