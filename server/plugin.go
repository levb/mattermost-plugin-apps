// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package main

import (
	"fmt"
	gohttp "net/http"
	"net/http/httputil"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	pluginapi "github.com/mattermost/mattermost-plugin-api"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/server/appservices"
	"github.com/mattermost/mattermost-plugin-apps/server/aws"
	"github.com/mattermost/mattermost-plugin-apps/server/builtin"
	"github.com/mattermost/mattermost-plugin-apps/server/config"
	"github.com/mattermost/mattermost-plugin-apps/server/examples/go/hello/builtin_hello"
	"github.com/mattermost/mattermost-plugin-apps/server/examples/go/hello/http_hello"
	"github.com/mattermost/mattermost-plugin-apps/server/examples/js/aws_hello"
	"github.com/mattermost/mattermost-plugin-apps/server/http"
	"github.com/mattermost/mattermost-plugin-apps/server/http/oauth"
	"github.com/mattermost/mattermost-plugin-apps/server/http/restapi"
	"github.com/mattermost/mattermost-plugin-apps/server/proxy"
	"github.com/mattermost/mattermost-plugin-apps/server/store"
)

// <><> const mutexKey = "Cluster_Mutex"

type Plugin struct {
	plugin.MattermostPlugin
	*config.BuildConfig

	mm          *pluginapi.Client
	conf        config.Service
	proxy       proxy.Service
	appServices appservices.Service
	aws         aws.Service
	http        http.Service
	store       *store.Service
}

func NewPlugin(buildConfig *config.BuildConfig) *Plugin {
	return &Plugin{
		BuildConfig: buildConfig,
	}
}

func (p *Plugin) OnActivate() error {
	mm := pluginapi.NewClient(p.API)
	p.mm = mm

	botUserID, err := mm.Bot.EnsureBot(&model.Bot{
		Username:    config.BotUsername,
		DisplayName: config.BotDisplayName,
		Description: config.BotDescription,
	}, pluginapi.ProfileImagePath("assets/profile.png"))
	if err != nil {
		return errors.Wrap(err, "failed to ensure bot account")
	}

	p.conf = config.NewService(mm, p.BuildConfig, botUserID)
	p.aws = aws.NewService(&mm.Log)
	p.store = store.NewService(p.mm, p.conf)

	// OnConfigurationChange updates aws and store with config values
	err = p.OnConfigurationChange()
	if err != nil {
		return errors.Wrap(err, "failed to obtain plugin configuration")
	}

	p.proxy = proxy.NewService(p.mm, p.aws, p.conf, p.store)
	p.appServices = appservices.NewService(mm, p.conf, p.store)
	p.http = http.NewService(mux.NewRouter(), p.mm, p.conf, p.proxy, p.appServices,
		oauth.Init,
		restapi.Init,
		http_hello.Init,
	)

	// Initialize the list of available apps (manifests).
	err = p.store.Manifest.InitGlobal(p.aws.Client())
	if err != nil {
		return errors.Wrap(err, "failed to initialize data store")
	}
	p.store.Manifest.InitBuiltin(
		aws_hello.Manifest(),
		// builtin - not in the list, to hide it from `apps list`
		builtin_hello.Manifest(),
		// http_hello - not in the list, can be `apps install developer SITEURL/plugins/com.mattermost.apps/example/hello/mattermost-app.json
	)

	// Initialize the list of installed apps.
	p.proxy.InitBuiltinApps(
		builtin_hello.NewHelloApp(),
		builtin.NewBuiltinApp(p.mm, p.conf, p.proxy, p.store),
	)

	return nil
}

func (p *Plugin) OnConfigurationChange() error {
	if p.conf == nil {
		// pre-activate, nothing to do.
		return nil
	}

	sc := config.StoredConfig{}
	_ = p.mm.Configuration.LoadPluginConfiguration(&sc)

	return p.conf.Refresh(&sc,
		p.aws,
		p.store.App,
		p.store.Manifest)
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w gohttp.ResponseWriter, req *gohttp.Request) {
	bb, _ := httputil.DumpRequest(req, false)
	fmt.Printf("<><> ServeHTTP: %s\n", string(bb))
	p.http.ServeHTTP(c, w, req)
}

func (p *Plugin) UserHasBeenCreated(pluginContext *plugin.Context, user *model.User) {
	_ = p.proxy.Notify(apps.NewUserContext(user), apps.SubjectUserCreated)
}

func (p *Plugin) UserHasJoinedChannel(pluginContext *plugin.Context, cm *model.ChannelMember, actingUser *model.User) {
	_ = p.proxy.Notify(apps.NewChannelMemberContext(cm, actingUser), apps.SubjectUserJoinedChannel)
}

func (p *Plugin) UserHasLeftChannel(pluginContext *plugin.Context, cm *model.ChannelMember, actingUser *model.User) {
	_ = p.proxy.Notify(apps.NewChannelMemberContext(cm, actingUser), apps.SubjectUserLeftChannel)
}

func (p *Plugin) UserHasJoinedTeam(pluginContext *plugin.Context, tm *model.TeamMember, actingUser *model.User) {
	_ = p.proxy.Notify(apps.NewTeamMemberContext(tm, actingUser), apps.SubjectUserJoinedTeam)
}

func (p *Plugin) UserHasLeftTeam(pluginContext *plugin.Context, tm *model.TeamMember, actingUser *model.User) {
	_ = p.proxy.Notify(apps.NewTeamMemberContext(tm, actingUser), apps.SubjectUserLeftTeam)
}

func (p *Plugin) MessageHasBeenPosted(pluginContext *plugin.Context, post *model.Post) {
	_ = p.proxy.Notify(apps.NewPostContext(post), apps.SubjectPostCreated)
}

func (p *Plugin) ChannelHasBeenCreated(pluginContext *plugin.Context, ch *model.Channel) {
	_ = p.proxy.Notify(apps.NewChannelContext(ch), apps.SubjectChannelCreated)
}
