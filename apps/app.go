package apps

import (
	"unicode"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/v5/model"
)

// AppID is a globally unique identifier that represents a Mattermost App.
// Allowed characters are letters, numbers, underscores and hyphens.
type AppID string

func (id AppID) IsValid() error {
	for _, c := range id {
		if unicode.IsLetter(c) {
			continue
		}

		if unicode.IsNumber(c) {
			continue
		}

		if c == '-' || c == '_' {
			continue
		}

		return errors.Errorf("invalid character %v in appID", c)
	}

	return nil
}

// AppVersion is the version of a Mattermost App.
// Allowed characters are letters, numbers, underscores and hyphens.
type AppVersion string

func (v AppVersion) IsValid() error {
	for _, c := range v {
		if unicode.IsLetter(c) {
			continue
		}

		if unicode.IsNumber(c) {
			continue
		}

		if c == '-' || c == '_' {
			continue
		}

		return errors.Errorf("invalid character %v in appVersion", c)
	}

	return nil
}

type AppVersionMap map[AppID]AppVersion

type AppType string

// default is HTTP
const (
	AppTypeHTTP    = "http"
	AppTypeAWS     = "aws"
	AppTypeBuiltin = "builtin"
)

func (at AppType) IsValid() bool {
	return at == AppTypeHTTP ||
		at == AppTypeAWS ||
		at == AppTypeBuiltin
}

type Common struct {
	AppID   AppID      `json:"app_id"`
	Type    AppType    `json:"type"`
	Version AppVersion `json:"version"`

	DisplayName string `json:"display_name,omitempty"`
	Description string `json:"description,omitempty"`
	HomepageURL string `json:"homepage_url,omitempty"`

	OnDisable   *Call `json:"on_disable,omitempty"`
	OnEnable    *Call `json:"on_enable,omitempty"`
	OnInstall   *Call `json:"on_install,omitempty"`
	OnStartup   *Call `json:"on_startup,omitempty"`
	OnUninstall *Call `json:"on_uninstall,omitempty"`
	Bindings    *Call `json:"bindings,omitempty"`
}

// Manifest describes a "known", installable app. They generally come from the
// marketplace, and can also be installed as overrides by developers.
// Manifest should be abbreviated as `m`.
type Manifest struct {
	Common

	RequestedPermissions Permissions `json:"requested_permissions,omitempty"`

	// RequestedLocations is the list of top-level locations that the
	// application intends to bind to, e.g. `{"/post_menu", "/channel_header",
	// "/command/apptrigger"}``.
	RequestedLocations Locations `json:"requested_locations,omitempty"`

	// For HTTP Apps all paths are relative to the RootURL.
	RootURL string `json:"root_url,omitempty"`

	// For AWS Apps, we need mappings from Call and static paths to the
	// respective AWS resources: names for Lambda functions, and S3 key for
	// S3 static files (the bucket is configured, per instance).
	FunctionRoutes map[string]string `json:"function_routes,omitempty"`
	StaticRoutes   map[string]string `json:"static_routes,omitempty"`
}

// App describes an App installed (or about to be installed) on a Mattermost instance.
// App should be abbreviated as `app`.
type App struct {
	Common

	Disabled bool `json:"disabled"`

	// Secret is used to issue JWT
	Secret string `json:"secret,omitempty"`

	OAuth2ClientID     string `json:"oauth2_client_id,omitempty"`
	OAuth2ClientSecret string `json:"oauth2_client_secret,omitempty"`
	OAuth2TrustedApp   bool   `json:"oauth2_trusted_app,omitempty"`

	BotUserID      string `json:"bot_user_id,omitempty"`
	BotUsername    string `json:"bot_username,omitempty"`
	BotAccessToken string `json:"bot_access_token,omitempty"`

	// Grants should be scopable in the future, per team, channel, post with
	// regexp.
	GrantedPermissions Permissions `json:"granted_permissions,omitempty"`

	// GrantedLocations contains the list of top locations that the
	// application is allowed to bind to.
	GrantedLocations Locations `json:"granted_locations,omitempty"`
}

type MarketplaceApp struct {
	Manifest  *Manifest                `json:"manifest"`
	Installed bool                     `json:"installed"`
	Enabled   bool                     `json:"enabled"`
	Labels    []model.MarketplaceLabel `json:"labels,omitempty"`
}

// Conventions for Apps paths, and field names
const (
	DefaultInstallCallPath  = "/install"
	DefaultBindingsCallPath = "/bindings"
)

var DefaultInstallCall = &Call{
	URL: DefaultInstallCallPath,
	Expand: &Expand{
		App:              ExpandAll,
		AdminAccessToken: ExpandAll,
	},
}

var DefaultBindingsCall = &Call{
	URL: DefaultBindingsCallPath,
}
