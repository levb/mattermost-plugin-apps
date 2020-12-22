// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package api

// Internal configuration apps.of mattermost-plugin-apps
const (
	Repository     = "mattermost-plugin-apps"
	CommandTrigger = "apps"

	BotUsername    = "appsbot"
	BotDisplayName = "Mattermost Apps"
	BotDescription = "Mattermost Apps Registry and API proxy."

	// TODO replace Interactive Dialogs with Modal, eliminate the need for
	// /dialog endpoints.
	InteractiveDialogPath = "/dialog"

	// Top-level path(s) for HTTP example apps.
	HelloHTTPPath = "/example/hello"

	// Top-level path for the REST APIs exposed by the plugin itself.
	APIPath = "/api/v1"

	// Top-level path for the Apps namespaces, followed by /AppID/subpath.
	AppsPath = "/apps"

	// OAuth2 sub-paths.
	OAuth2Path = "/oauth2" // convention for Mattermost Apps, comes from OAuther

	// Other sub-paths.
	CallPath      = "/call"
	KVPath        = "/kv"
	SubscribePath = "/subscribe"
	BindingsPath  = "/bindings"
)

const (
	PropTeamID             = "team_id"
	PropChannelID          = "channel_id"
	PropPostID             = "post_id"
	PropOAuth2ClientSecret = "oauth2_client_secret" // nolint:gosec
	PropAppBindings        = "app_bindings"
)
