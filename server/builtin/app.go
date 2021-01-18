package builtin

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/mattermost/mattermost-plugin-apps/server/api"
	"github.com/mattermost/mattermost-plugin-apps/server/utils/md"
	"github.com/pkg/errors"
)

const (
	AppID          = "apps"
	AppDisplayName = "Mattermost Apps plugin"
	AppDescription = "Install and manage Mattermost Apps"
)

const (
	CommandBindings   = "bindings"
	CommandClean      = "clean"
	CommandConnect    = "connect"
	CommandDebug      = "debug"
	CommandDisconnect = "disconnect"
	CommandList       = "list"
	CommandInfo       = "info"
	CommandInstall    = "install"
)

const (
	fieldAppID              = "app_id"
	fieldExampleApp         = "example"
	fieldManifestURL        = "manifest_url"
	fieldRequireUserConsent = "require_consent"
)

const (
	flagUserID             = "user"
	flagAppID              = "app"
	flagExampleApp         = "example"
	flagManifestURL        = "manifest"
	flagRequireUserConsent = "require-consent"
)

const (
	PathConnect       = "/connect"
	PathDebugBindings = "/debug-bindings"
	PathDebugClean    = "/debug-clean"
	PathDebugInstall  = "/debug-install"
	PathInstallApp    = "/install-app"
	PathDisconnect    = "/disconnect"
	PathInfo          = "/info"
	PathList          = "/list"
)

type App struct {
	API *api.Service
}

func NewApp(api *api.Service) *App {
	return &App{api}
}

var _ api.Upstream = (*App)(nil)

func (a *App) MattermostApp() *api.App {
	conf := a.API.Configurator.GetConfig()
	return &api.App{
		Manifest: &api.Manifest{
			AppID:       AppID,
			Type:        api.AppTypeBuiltin,
			DisplayName: AppDisplayName,
			Description: AppDescription,
			RequestedLocations: api.Locations{
				api.LocationCommand,
			},
		},
		BotUserID:   conf.BotUserID,
		BotUsername: api.BotUsername,
		GrantedLocations: api.Locations{
			api.LocationCommand,
		},
	}
}

func (a *App) Roundtrip(c *api.Call) (io.ReadCloser, error) {
	cr := &api.CallResponse{}
	switch c.URL {
	case api.BindingsPath:
		cr = a.funcGetBindings(c)
	case PathInfo:
		cr = simpleFunc(a.infoForm, a.info)(c)
	case PathList:
		cr = simpleFunc(a.listForm, a.list)(c)
	case PathConnect:
		cr = simpleFunc(a.connectForm, a.connect)(c)
	case PathDisconnect:
		cr = simpleFunc(a.connectForm, a.disconnect)(c)
	case PathDebugClean:
		cr = simpleFunc(nil, a.clean)(c)
	case PathDebugInstall:
		cr = simpleFunc(a.debugInstallForm, a.debugInstall)(c)
	default:
		return nil, errors.Errorf("%s is not found", c.URL)
	}

	bb, err := json.Marshal(cr)
	if err != nil {
		return nil, err
	}
	return ioutil.NopCloser(bytes.NewReader(bb)), nil
}

func (a *App) OneWay(call *api.Call) error {
	return nil
}

func simpleFunc(formf func(*api.Call) (*api.Form, error),
	submitf func(*api.Call) (md.MD, error)) func(call *api.Call) *api.CallResponse {
	return func(call *api.Call) *api.CallResponse {
		switch call.Type {
		case api.CallTypeForm:
			form := &api.Form{}
			if formf != nil {
				var err error
				form, err = formf(call)
				if err != nil {
					return api.NewErrorCallResponse(err)
				}
			}
			return &api.CallResponse{
				Type: api.CallResponseTypeForm,
				Form: form,
			}

		case api.CallTypeSubmit:
			txt, err := submitf(call)
			return api.NewCallResponse(txt, nil, err)

		default:
			return api.NewErrorCallResponse(errors.New("not supported"))
		}
	}
}
