// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package config

const ManifestsFile = "assets/manifests.json"

// Internal configuration apps.of mattermost-plugin-apps.
const (
	Repository = "mattermost-plugin-apps"

	BotUsername    = "appsbot"
	BotDisplayName = "Mattermost Apps"
	BotDescription = "Mattermost Apps Registry and API proxy."
)

// HTTP paths.
const (
	// Top-level path(s) for HTTP example apps.
	HelloHTTPPath = "/example/hello"

	// Top-level path for the REST APIs exposed by the plugin itself.
	APIPath = "/api/v1"

	// API sub-paths.
	CallPath        = "/call"
	KVPath          = "/kv"
	SubscribePath   = "/subscribe"
	UnsubscribePath = "/unsubscribe"
	BindingsPath    = "/bindings"

	// Top-level path for the Apps namespaces, followed by /AppID/subpath.
	// Used for webhooks, and to get static resources
	AppsPath = "/apps"

	// OAuth2 sub-paths.
	OAuth2Path = "/oauth2" // convention for Mattermost Apps, comes from OAuther

	// Marketplace sub-paths.
	PathMarketplace = "/marketplace"
)

const (
	PropTeamID      = "team_id"
	PropChannelID   = "channel_id"
	PropPostID      = "post_id"
	PropAppBindings = "app_bindings"
)
