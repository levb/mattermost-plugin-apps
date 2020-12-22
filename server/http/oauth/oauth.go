package oauth

import (
	"github.com/gorilla/mux"

	"github.com/mattermost/mattermost-plugin-apps/server/api"
)

type oauth struct {
	api *api.Service
}

func Init(router *mux.Router, appsService *api.Service) {
	a := &oauth{
		api: appsService,
	}

	router.HandleFunc(api.OAuth2Path, a.api.Proxy.HandleOAuth).Methods("GET")
}
