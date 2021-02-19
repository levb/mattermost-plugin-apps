package builtin

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"

	pluginapi "github.com/mattermost/mattermost-plugin-api"
	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/server/config"
	"github.com/mattermost/mattermost-plugin-apps/server/proxy"
	"github.com/mattermost/mattermost-plugin-apps/server/store"
	"github.com/mattermost/mattermost-plugin-apps/server/upstream"
	"github.com/mattermost/mattermost-plugin-apps/server/utils"
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
	fieldSecret             = "secret"
)

const (
	flagAppID              = "app"
	flagExampleApp         = "example"
	flagManifestURL        = "manifest"
	flagRequireUserConsent = "require-consent"
	flagSecret             = "secret"
	flagUserID             = "user"
)

const (
	PathConnect           = "/connect"
	PathDebugBindings     = "/debug-bindings"
	PathDebugClean        = "/debug-clean"
	PathInstallApp        = "/install-app"
	PathInstallAppCommand = "/install-app-command"
	PathDisconnect        = "/disconnect"
	PathInfo              = "/info"
	PathList              = "/list"
)

type App struct {
	conf  config.Service
	mm    *pluginapi.Client
	proxy proxy.Service
	store *store.Service
}

func NewApp(mm *pluginapi.Client, conf config.Service, proxy proxy.Service, store *store.Service) *App {
	return &App{
		conf:  conf,
		mm:    mm,
		proxy: proxy,
		store: store,
	}
}

var _ upstream.Upstream = (*App)(nil)

func (a *App) App() *apps.App {
	conf := a.conf.Get()
	return &apps.App{
		Common: apps.Common{
			AppID:       AppID,
			Type:        apps.AppTypeBuiltin,
			DisplayName: AppDisplayName,
			Description: AppDescription,
		},
		BotUserID:   conf.BotUserID,
		BotUsername: config.BotUsername,
		GrantedLocations: apps.Locations{
			apps.LocationCommand,
		},
	}
}

func (a *App) Roundtrip(c *apps.Call) (io.ReadCloser, error) {
	cr := &apps.CallResponse{}
	switch c.URL {
	case config.BindingsPath:
		cr = a.funcGetBindings(c)

	case PathInfo:
		cr = simpleFunc(a.infoForm, a.info)(c)
	case PathList:
		cr = simpleFunc(a.listForm, a.list)(c)
	case PathDebugClean:
		cr = simpleFunc(nil, a.clean)(c)

	case PathConnect:
		cr = simpleFunc(a.connectForm, a.connect)(c)
	case PathDisconnect:
		cr = simpleFunc(a.connectForm, a.disconnect)(c)

	case PathInstallAppCommand:
		cr = simpleFunc(a.installAppCommandForm, a.installAppCommand)(c)
	case PathInstallApp:
		cr = simpleFunc(nil, a.installApp)(c)

	default:
		return nil, errors.Errorf("%s is not found", c.URL)
	}

	bb, err := json.Marshal(cr)
	if err != nil {
		return nil, err
	}
	return ioutil.NopCloser(bytes.NewReader(bb)), nil
}

func (a *App) OneWay(call *apps.Call) error {
	return nil
}

func (a *App) GetStatic(path string) ([]byte, error) {
	return nil, utils.ErrNotFound
}

func simpleFunc(
	formf func(*apps.Call) (*apps.Form, error),
	submitf func(*apps.Call) *apps.CallResponse) func(call *apps.Call) *apps.CallResponse {
	return func(call *apps.Call) *apps.CallResponse {
		switch call.Type {
		case apps.CallTypeForm:
			form := &apps.Form{}
			if formf != nil {
				var err error
				form, err = formf(call)
				if err != nil {
					return apps.NewErrorCallResponse(err)
				}
			}
			return &apps.CallResponse{
				Type: apps.CallResponseTypeForm,
				Form: form,
			}

		case apps.CallTypeSubmit:
			return submitf(call)

		default:
			return apps.NewErrorCallResponse(errors.New("not supported"))
		}
	}
}
