package restapi

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	pluginapi "github.com/mattermost/mattermost-plugin-api"
	"github.com/mattermost/mattermost-plugin-apps/server/appservices"
	"github.com/mattermost/mattermost-plugin-apps/server/config"
	"github.com/mattermost/mattermost-plugin-apps/server/proxy"
	"github.com/mattermost/mattermost-plugin-apps/server/utils/httputils"
)

type restapi struct {
	mm          *pluginapi.Client
	appservices appservices.Service
	proxy       proxy.Service
	conf        config.Service
}

func Init(router *mux.Router, mm *pluginapi.Client, conf config.Service, proxy proxy.Service, appservices appservices.Service) {
	a := &restapi{
		mm:          mm,
		conf:        conf,
		appservices: appservices,
		proxy:       proxy,
	}

	subrouter := router.PathPrefix(config.APIPath).Subrouter()
	subrouter.HandleFunc(config.BindingsPath, checkAuthorized(a.handleGetBindings)).Methods("GET")
	subrouter.HandleFunc(config.CallPath, a.handleCall).Methods("POST")
	subrouter.HandleFunc(config.SubscribePath, a.handleSubscribe).Methods("POST")
	subrouter.HandleFunc(config.UnsubscribePath, a.handleUnsubscribe).Methods("POST")

	subrouter.HandleFunc(config.KVPath+"/{key}", a.handleKV(a.kvGet)).Methods("GET")
	subrouter.HandleFunc(config.KVPath+"/{key}", a.handleKV(a.kvPut)).Methods("PUT", "POST")
	subrouter.HandleFunc(config.KVPath+"/", a.handleKV(a.kvList)).Methods("GET")
	subrouter.HandleFunc(config.KVPath+"/{key}", a.handleKV(a.kvHead)).Methods("HEAD")
	subrouter.HandleFunc(config.KVPath+"/{key}", a.handleKV(a.kvDelete)).Methods("DELETE")

	// subrouter.HandleFunc(config.PathMarketplace, checkAuthorized(a.handleGetMarketplace)).Methods(http.MethodGet)
}

func checkAuthorized(f func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		actingUserID := req.Header.Get("Mattermost-User-Id")
		if actingUserID == "" {
			httputils.WriteUnauthorizedError(w, errors.New("not authorized"))
			return
		}

		f(w, req, actingUserID)
	}
}
