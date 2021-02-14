package api

import (
	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/server/utils/md"
)

type Admin interface {
	ListApps() ([]*apps.App, error)
	InstallApp(*apps.Context, *apps.InInstallApp) (*apps.App, md.MD, error)
	LoadAppsList() error
}
