package hello

import (
	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/pkg/errors"
)

func (h *HelloApp) Survey(call *apps.Call) *apps.CallResponse {
	switch call.Type {
	case apps.CallTypeForm:
		return newSurveyFormResponse(call)

	case apps.CallTypeSubmit:
		err := h.processSurvey(call)
		return apps.NewCallResponse("ok", nil, err)

	default:
		return apps.NewErrorCallResponse(errors.New("not supported"))
	}
}

func NewSurveyForm(message string) *apps.Form {
	return &apps.Form{
		Title:         "Emotional response survey",
		Header:        message,
		Footer:        "Let the world know!",
		SubmitButtons: fieldResponse,
		Fields: []*apps.Field{
			{
				Name: fieldResponse,
				Type: apps.FieldTypeStaticSelect,
				SelectStaticOptions: []apps.SelectOption{
					{Label: "Like", Value: "like"},
					{Label: "Dislike", Value: "dislike"},
				},
			},
		},
	}
}

func newSurveyFormResponse(c *apps.Call) *apps.CallResponse {
	message := c.GetStringValue(fieldMessage, "default hello message")
	return &apps.CallResponse{
		Type: apps.CallResponseTypeForm,
		Form: NewSurveyForm(message),
	}
}

func (h *HelloApp) processSurvey(c *apps.Call) error {
	// TODO post something; for embedded form - what do we do?
	return nil
}
