package hello

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/apps/mmclient"
	"github.com/mattermost/mattermost-plugin-apps/server/config"
	"github.com/mattermost/mattermost-plugin-apps/server/utils/md"
)

func (h *HelloApp) SendSurvey(call *apps.Call) *apps.CallResponse {
	switch call.Type {
	case apps.CallTypeForm:
		return NewSendSurveyFormResponse(call)

	case apps.CallTypeSubmit:
		txt, err := h.sendSurvey(call)
		return apps.NewCallResponse(txt, nil, err)

	case apps.CallTypeLookup:
		return &apps.CallResponse{
			Data: map[string]interface{}{
				"items": []*apps.SelectOption{
					{
						Label: "Option 1",
						Value: "option1",
					},
				},
			},
		}

	default:
		return apps.NewErrorCallResponse(errors.New("not supported"))
	}
}

func (h *HelloApp) SendSurveyModal(c *apps.Call) *apps.CallResponse {
	return NewSendSurveyFormResponse(c)
}

func (h *HelloApp) SendSurveyCommandToModal(c *apps.Call) *apps.CallResponse {
	return NewSendSurveyPartialFormResponse(c)
}

type SurveyFormSubmission struct {
	UserID  string                 `json:"userID"`
	Message string                 `json:"message"`
	Other   map[string]interface{} `json:"other"`
}

func extractSurveyFormValues(c *apps.Call) SurveyFormSubmission {
	message := ""
	userID := ""
	var other map[string]interface{} = nil
	if c.Context != nil && c.Context.Post != nil {
		message = c.Context.Post.Message
	}

	topValues := c.Values
	formValues := c.Values
	if c.Type == apps.CallTypeForm && topValues != nil {
		formValues, _ = topValues["values"].(map[string]interface{})
	}

	if formValues != nil {
		userID, _ = formValues["userID"].(string)
		message, _ = formValues["message"].(string)
		otherTemp, ok2 := formValues["other"].(map[string]interface{})
		if ok2 {
			other = otherTemp
		} else {
			other = nil
		}
	}

	return SurveyFormSubmission{
		UserID:  userID,
		Message: message,
		Other:   other,
	}
}

func NewSendSurveyFormResponse(c *apps.Call) *apps.CallResponse {
	submission := extractSurveyFormValues(c)
	name, _ := c.Values["name"].(string)

	if name == "userID" {
		submission.Message = fmt.Sprintf("%s Now sending to %s.", submission.Message, submission.UserID)
	}

	return &apps.CallResponse{
		Type: apps.CallResponseTypeForm,
		Form: &apps.Form{
			Title:  "Send a survey to user",
			Header: "Message modal form header",
			Footer: "Message modal form footer",
			Call:   apps.MakeCall(PathSendSurvey),
			Fields: []*apps.Field{
				{
					Name:                 fieldUserID,
					Type:                 apps.FieldTypeUser,
					Description:          "User to send the survey to",
					Label:                "user",
					ModalLabel:           "User",
					AutocompleteHint:     "enter user ID or @user",
					AutocompletePosition: 1,
					Value:                submission.UserID,
					SelectRefresh:        true,
				}, {
					Name:             "other",
					Type:             apps.FieldTypeDynamicSelect,
					Description:      "Some values",
					Label:            "other",
					AutocompleteHint: "Pick one",
					ModalLabel:       "Other",
					Value:            submission.Other,
				}, {
					Name:             fieldMessage,
					Type:             apps.FieldTypeText,
					Description:      "Text to ask the user about",
					IsRequired:       true,
					Label:            "message",
					ModalLabel:       "Text",
					AutocompleteHint: "Anything you want to say",
					TextSubtype:      "textarea",
					TextMinLength:    2,
					TextMaxLength:    1024,
					Value:            submission.Message,
				},
			},
		},
	}
}

func NewSendSurveyPartialFormResponse(c *apps.Call) *apps.CallResponse {
	if c.Type == apps.CallTypeSubmit {
		return NewSendSurveyFormResponse(c)
	}

	return &apps.CallResponse{
		Type: apps.CallResponseTypeForm,
		Form: &apps.Form{
			Title:  "Send a survey to user",
			Header: "Message modal form header",
			Footer: "Message modal form footer",
			Call:   apps.MakeCall(PathSendSurveyCommandToModal),
			Fields: []*apps.Field{
				{
					Name:             fieldMessage,
					Type:             apps.FieldTypeText,
					Description:      "Text to ask the user about",
					IsRequired:       true,
					Label:            "message",
					ModalLabel:       "Text",
					AutocompleteHint: "Anything you want to say",
					TextSubtype:      "textarea",
					TextMinLength:    2,
					TextMaxLength:    1024,
					Value:            "",
				},
			},
		},
	}
}

func (h *HelloApp) sendSurvey(c *apps.Call) (md.MD, error) {
	bot := mmclient.AsBot(c.Context)
	userID := c.GetStringValue(fieldUserID, c.Context.ActingUserID)

	// TODO this should be done with expanding mentions, make a ticket
	if strings.HasPrefix(userID, "@") {
		user, _ := bot.GetUserByUsername(userID[1:], "")
		if user != nil {
			userID = user.Id
		}
	}

	message := c.GetStringValue(fieldMessage, "Hello")
	if c.Context.Post != nil {
		message += "\n>>> " + c.Context.Post.Message
	}

	err := sendSurvey(bot, userID, message)
	if err != nil {
		return "", err
	}

	return "Successfully sent survey", nil
}

func sendSurvey(bot *mmclient.Client, userID, message string) error {
	p := &model.Post{
		Message: "Please respond to this survey: " + message,
	}
	p.AddProp(config.PropAppBindings, []*apps.Binding{
		{
			Location: "survey",
			Form:     NewSurveyForm(message),
		},
	})
	_, err := bot.DMPost(userID, p)
	return err
}
