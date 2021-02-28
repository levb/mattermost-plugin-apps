package hello

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/apps/mmclient"
	"github.com/mattermost/mattermost-plugin-apps/server/utils/md"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
)

func (h *HelloApp) PostAsUser(call *apps.Call) *apps.CallResponse {
	switch call.Type {
	case apps.CallTypeForm:
		return newPostAsUserFormResponse(call)

	case apps.CallTypeSubmit:
		txt, err := h.postAsUser(call)
		return apps.NewCallResponse(txt, nil, err)

	default:
		return apps.NewErrorCallResponse(errors.New("not supported"))
	}
}

func newPostAsUserFormResponse(c *apps.Call) *apps.CallResponse {
	message := ""
	if c.Context != nil && c.Context.Post != nil {
		message = c.Context.Post.Message
	}

	return &apps.CallResponse{
		Type: apps.CallResponseTypeForm,
		Form: &apps.Form{
			Title:  fmt.Sprintf("Post to the %s channel, as user", c.Context.AppID),
			Header: "Message modal form header",
			Footer: "Message modal form footer",
			Fields: []*apps.Field{
				{
					Name:             fieldMessage,
					Type:             apps.FieldTypeText,
					Description:      "Text to post",
					IsRequired:       true,
					Label:            "message",
					ModalLabel:       "Text",
					AutocompleteHint: "Anything you want to say",
					TextSubtype:      "textarea",
					TextMinLength:    2,
					TextMaxLength:    1024,
					Value:            message,
				},
			},
		},
	}
}

func (h *HelloApp) postAsUser(call *apps.Call) (md.MD, error) {
	acting := mmclient.AsActingUser(call.Context)

	channel, _ := acting.GetChannelByName(string(call.Context.AppID), string(call.Context.TeamID), "")
	if channel == nil {
		return "", errors.Errorf("channel %s not found", call.Context.AppID)
	}

	_, err := acting.CreatePost(&model.Post{
		Message: "Hallo!",
	})
	if err != nil {
		return "", err
	}

	return "OK", nil
}
