package hello

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-apps/server/api"
	"github.com/mattermost/mattermost-plugin-apps/server/examples"
	"github.com/mattermost/mattermost-plugin-apps/server/utils/md"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
)

func (h *HelloApp) PostAsUser(call *api.Call) *api.CallResponse {
	switch call.Type {
	case api.CallTypeForm:
		return newPostAsUserFormResponse(call)

	case api.CallTypeSubmit:
		txt, err := h.postAsUser(call)
		return api.NewCallResponse(txt, nil, err)

	default:
		return api.NewErrorCallResponse(errors.New("not supported"))
	}
}

func newPostAsUserFormResponse(c *api.Call) *api.CallResponse {
	message := ""
	if c.Context != nil && c.Context.Post != nil {
		message = c.Context.Post.Message
	}

	return &api.CallResponse{
		Type: api.CallResponseTypeForm,
		Form: &api.Form{
			Title:  fmt.Sprintf("Post to the %s channel, as user", c.Context.AppID),
			Header: "Message modal form header",
			Footer: "Message modal form footer",
			Fields: []*api.Field{
				{
					Name:             fieldMessage,
					Type:             api.FieldTypeText,
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

func (h *HelloApp) postAsUser(call *api.Call) (md.MD, error) {
	acting := examples.AsActingUser(call.Context)

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
