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
	"github.com/mattermost/mattermost-plugin-apps/server/utils"
	"github.com/pkg/errors"
)

const (
	AppID          = "apps"
	AppDisplayName = "Mattermost Apps plugin"
	AppDescription = "Install and manage Mattermost Apps"
)

const (
	CommandBindings    = "bindings"
	CommandClean       = "clean"
	CommandConnect     = "connect"
	CommandDebug       = "debug"
	CommandDisconnect  = "disconnect"
	CommandList        = "list"
	CommandInfo        = "info"
	CommandInstall     = "install"
	CommandMarketplace = "marketplace"
	CommandDeveloper   = "developer"
)

const (
	contextInstallAppID = "install_app_id"

	fURL                = "url"
	fConsentPermissions = "consent_permissions"
	fConsentLocations   = "consent_locations"
	fRequireUserConsent = "require_user_consent"
	fSecret             = "secret"
	fAppID              = "app"
	fUserID             = "user"
)

const (
	PathConnect            = "/connect"
	PathDebugBindings      = "/debug-bindings"
	PathDebugClean         = "/debug-clean"
	PathInstallMarketplace = "/install-marketplace"
	PathInstallDeveloper   = "/install-developer"
	PathInstallApp         = "/install-app"
	PathDisconnect         = "/disconnect"
	PathInfo               = "/info"
	PathList               = "/list"
)

type builtinApp struct {
	conf  config.Service
	mm    *pluginapi.Client
	proxy proxy.Service
	store *store.Service
}

func NewBuiltinApp(mm *pluginapi.Client, conf config.Service, proxy proxy.Service, store *store.Service) *builtinApp {
	return &builtinApp{
		conf:  conf,
		mm:    mm,
		proxy: proxy,
		store: store,
	}
}

func (a *builtinApp) App() *apps.App {
	conf := a.conf.Get()
	return &apps.App{
		Common: apps.Common{
			AppID:       AppID,
			Type:        apps.AppTypeBuiltin,
			DisplayName: AppDisplayName,
			Description: AppDescription,
			Version:     apps.AppVersion(conf.BuildConfig.BuildHashShort),
		},
		BotUserID:   conf.BotUserID,
		BotUsername: config.BotUsername,
		GrantedLocations: apps.Locations{
			apps.LocationCommand,
		},
	}
}

func (a *builtinApp) Roundtrip(c *apps.Call) (io.ReadCloser, error) {
	cr := &apps.CallResponse{}
	switch c.URL {
	case config.BindingsPath:
		cr = a.funcGetBindings(c)

	case PathInfo:
		cr = handle(a.infoForm, a.info, nil)(c)
	case PathList:
		cr = handle(a.listForm, a.list, nil)(c)
	case PathDebugClean:
		cr = handle(nil, a.clean, nil)(c)

	// case PathConnect:
	// 	cr = simpleFunc(a.connectForm, a.connect, nil)(c)
	// case PathDisconnect:
	// 	cr = simpleFunc(a.connectForm, a.disconnect, nil)(c)

	case PathInstallMarketplace:
		cr = handle(a.installMarketplaceCommandForm, a.installMarketplaceCommand, a.installMarketplaceLookup)(c)
	case PathInstallDeveloper:
		cr = handle(a.installDeveloperCommandForm, a.installDeveloperCommand, nil)(c)
	case PathInstallApp:
		cr = handle(a.installAppForm, a.installApp, nil)(c)

	default:
		return nil, errors.Errorf("%s is not found", c.URL)
	}

	bb, err := json.Marshal(cr)
	if err != nil {
		return nil, err
	}
	return ioutil.NopCloser(bytes.NewReader(bb)), nil
}

func (a *builtinApp) OneWay(call *apps.Call) error {
	return nil
}

func (a *builtinApp) GetStatic(path string) ([]byte, error) {
	return nil, utils.ErrNotFound
}

func handle(
	formf func(*apps.Call) *apps.CallResponse,
	submitf func(*apps.Call) *apps.CallResponse,
	lookupf func(*apps.Call) []*apps.SelectOption) func(call *apps.Call) *apps.CallResponse {
	return func(call *apps.Call) *apps.CallResponse {
		switch call.Type {
		case apps.CallTypeForm:
			if formf != nil {
				return formf(call)
			}
			return &apps.CallResponse{
				Type: apps.CallResponseTypeForm,
				Form: &apps.Form{},
			}

		case apps.CallTypeSubmit:
			if submitf != nil {
				return submitf(call)
			}

		case apps.CallTypeLookup:
			resp := &apps.CallResponse{}
			if lookupf != nil {
				options := lookupf(call)
				if len(options) != 0 {
					resp.Data = map[string]interface{}{
						"items": options,
					}
				}
			}
			return resp
		}
		return apps.NewErrorCallResponse(errors.New("not supported"))
	}
}
