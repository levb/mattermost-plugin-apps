package builtin_hello

import (
	"github.com/mattermost/mattermost-plugin-apps/server/api"
	"github.com/mattermost/mattermost-plugin-apps/server/examples/go/hello"
	"github.com/pkg/errors"
)

const (
	AppID          = "builtin"
	AppDisplayName = "builtin hello display name"
	AppDescription = "builtin hello description"
)

type helloapp struct {
	*hello.HelloApp
}

var _ api.Upstream = (*helloapp)(nil)

func New(appService *api.Service) *helloapp {
	return &helloapp{
		HelloApp: &hello.HelloApp{
			API: appService,
		},
	}
}

func Manifest() *api.Manifest {
	return &api.Manifest{
		AppID:       AppID,
		Type:        api.AppTypeBuiltin,
		DisplayName: AppDisplayName,
		Description: AppDescription,
		RequestedPermissions: api.Permissions{
			api.PermissionUserJoinedChannelNotification,
			api.PermissionActAsUser,
			api.PermissionActAsBot,
		},
		RequestedLocations: api.Locations{
			api.LocationChannelHeader,
			api.LocationPostMenu,
			api.LocationCommand,
			api.LocationInPost,
		},
		HomepageURL: ("https://github.com/mattermost"),
	}
}

func (h *helloapp) Call(c *api.Call) *api.CallResponse {
	switch c.URL {
	case api.DefaultInstallCallPath:
		return h.Install(c)
	case hello.PathSendSurvey:
		return h.SendSurvey(c)
	case hello.PathSurvey:
		return h.Survey(c)
	default:
		return api.NewErrorCallResponse(errors.Errorf("%s is not found", c.URL))
	}
}

func (h *helloapp) Notify(call *api.Call) error {
	switch call.Context.Subject {
	case api.SubjectUserJoinedChannel:
		h.HelloApp.UserJoinedChannel(call)
	default:
		return errors.Errorf("%s is not supported", call.Context.Subject)
	}
	return nil
}

func (h *helloapp) Install(c *api.Call) *api.CallResponse {
	if c.Type != api.CallTypeSubmit {
		return api.NewErrorCallResponse(errors.New("not supported"))
	}
	out, err := h.HelloApp.Install(AppID, AppDisplayName, c)
	if err != nil {
		return api.NewErrorCallResponse(err)
	}
	return &api.CallResponse{
		Type:     api.CallResponseTypeOK,
		Markdown: out,
	}
}

func (h *helloapp) GetBindings(c *api.Call) ([]*api.Binding, error) {
	return hello.Bindings(), nil
}

func (h *helloapp) SendSurvey(c *api.Call) *api.CallResponse {
	switch c.Type {
	case api.CallTypeForm:
		return hello.NewSendSurveyFormResponse(c)

	case api.CallTypeSubmit:
		txt, err := h.HelloApp.SendSurvey(c)
		if err != nil {
			return api.NewErrorCallResponse(err)
		}
		return &api.CallResponse{
			Type:     api.CallResponseTypeOK,
			Markdown: txt,
		}
	}

	return nil
}

func (h *helloapp) Survey(c *api.Call) *api.CallResponse {
	switch c.Type {
	case api.CallTypeForm:
		return hello.NewSurveyFormResponse(c)

	case api.CallTypeSubmit:
		err := h.ProcessSurvey(c)
		if err != nil {
			return api.NewErrorCallResponse(err)
		}
		return &api.CallResponse{
			Type:     api.CallResponseTypeOK,
			Markdown: "<><> TODO",
		}
	}
	return nil
}
