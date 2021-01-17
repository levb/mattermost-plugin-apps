package builtin_hello

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"

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

func (h *helloapp) Roundtrip(c *api.Call) (io.ReadCloser, error) {
	cr := &api.CallResponse{}
	switch c.URL {
	case api.BindingsPath:
		cr = h.GetBindings(c)
	case api.DefaultInstallCallPath:
		cr = h.Install(AppID, AppDisplayName, c)
	case hello.PathSendSurvey:
		cr = h.SendSurvey(c)
	case hello.PathSendSurveyModal:
		cr = h.SendSurveyModal(c)
	case hello.PathSendSurveyCommandToModal:
		cr = h.SendSurveyCommandToModal(c)
	case hello.PathSurvey:
		cr = h.Survey(c)
	case hello.PathPostAsUser:
		cr = h.PostAsUser(c)
	default:
		return nil, errors.Errorf("%s is not found", c.URL)
	}

	bb, err := json.Marshal(cr)
	if err != nil {
		return nil, err
	}
	return ioutil.NopCloser(bytes.NewReader(bb)), nil
}

func (h *helloapp) OneWay(call *api.Call) error {
	switch call.Context.Subject {
	case api.SubjectUserJoinedChannel:
		h.HelloApp.UserJoinedChannel(call)
	default:
		return errors.Errorf("%s is not supported", call.Context.Subject)
	}
	return nil
}

