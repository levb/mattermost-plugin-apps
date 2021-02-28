package builtin_hello

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/server/config"
	"github.com/mattermost/mattermost-plugin-apps/server/examples/go/hello"
	"github.com/mattermost/mattermost-plugin-apps/server/utils"
)

const (
	AppID          = "builtin"
	AppDisplayName = "builtin hello display name"
	AppDescription = "builtin hello description"
)

type helloapp struct {
	*hello.HelloApp
}

func NewHelloApp() *helloapp {
	return &helloapp{
		HelloApp: &hello.HelloApp{},
	}
}

var common = apps.Common{
	AppID:       AppID,
	Type:        apps.AppTypeBuiltin,
	Version:     "pre-release",
	DisplayName: AppDisplayName,
	Description: AppDescription,
	HomepageURL: ("https://github.com/mattermost"),
}

var permissions = apps.Permissions{
	apps.PermissionUserJoinedChannelNotification,
	apps.PermissionActAsUser,
	apps.PermissionActAsBot,
}

var locations = apps.Locations{
	apps.LocationChannelHeader,
	apps.LocationPostMenu,
	apps.LocationCommand,
	apps.LocationInPost,
}

func Manifest() *apps.Manifest {
	return &apps.Manifest{
		Common:               common,
		RequestedPermissions: permissions,
		RequestedLocations:   locations,
	}
}

func (h *helloapp) App() *apps.App {
	return &apps.App{
		Common:             common,
		GrantedPermissions: permissions,
		GrantedLocations:   locations,
	}
}

func (h *helloapp) Roundtrip(c *apps.Call) (io.ReadCloser, error) {
	cr := &apps.CallResponse{}
	switch c.URL {
	case config.BindingsPath:
		cr = h.GetBindings(c)
	case "/install":
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

func (h *helloapp) OneWay(call *apps.Call) error {
	switch call.Context.Subject {
	case apps.SubjectUserJoinedChannel:
		h.HelloApp.UserJoinedChannel(call)
	default:
		return errors.Errorf("%s is not supported", call.Context.Subject)
	}
	return nil
}

func (h *helloapp) GetStatic(path string) ([]byte, error) {
	return nil, utils.ErrNotFound
}
