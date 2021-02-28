package oauth

import (
	"github.com/gorilla/mux"

	pluginapi "github.com/mattermost/mattermost-plugin-api"

	"github.com/mattermost/mattermost-plugin-apps/server/appservices"
	"github.com/mattermost/mattermost-plugin-apps/server/config"
	"github.com/mattermost/mattermost-plugin-apps/server/proxy"
)

type oauth struct {
	proxy proxy.Service
}

func Init(router *mux.Router, _ *pluginapi.Client, _ config.Service, proxy proxy.Service, _ appservices.Service) {
	a := &oauth{
		proxy: proxy,
	}

	router.HandleFunc(config.OAuth2Path, a.proxy.HandleOAuth).Methods("GET")
}
