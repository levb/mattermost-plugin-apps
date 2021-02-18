package http_hello

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	pluginapi "github.com/mattermost/mattermost-plugin-api"
	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/server/appservices"
	"github.com/mattermost/mattermost-plugin-apps/server/config"
	"github.com/mattermost/mattermost-plugin-apps/server/examples/go/hello"
	"github.com/mattermost/mattermost-plugin-apps/server/proxy"
	"github.com/mattermost/mattermost-plugin-apps/server/utils/httputils"
)

const (
	AppID          = "http-hello"
	AppSecret      = "1234"
	AppDisplayName = "Hallo სამყარო (http)"
	AppDescription = "Hallo სამყარო HTTP test app"
)

const (
	PathManifest = "/mattermost-app.json"
)

type helloapp struct {
	*hello.HelloApp
}

func NewHelloApp(conf config.Service) *helloapp {
	return &helloapp{
		HelloApp: hello.NewHelloApp(conf),
	}
}

// Init hello app router
func Init(router *mux.Router, _ *pluginapi.Client, conf config.Service, _ proxy.Service, _ appservices.Service) {
	h := NewHelloApp(conf)

	r := router.PathPrefix(config.HelloHTTPPath).Subrouter()
	r.HandleFunc(PathManifest, h.handleManifest).Methods("GET")

	handle(r, "/install", h.Install)
	handle(r, "/bindings", h.GetBindings)
	handle(r, hello.PathSendSurvey, h.SendSurvey)
	handle(r, hello.PathSendSurveyModal, h.SendSurveyModal)
	handle(r, hello.PathSendSurveyCommandToModal, h.SendSurveyCommandToModal)
	handle(r, hello.PathSurvey, h.Survey)
	handle(r, hello.PathUserJoinedChannel, h.UserJoinedChannel)
	handle(r, hello.PathPostAsUser, h.PostAsUser)
}

func Manifest(conf config.Config) *apps.Manifest {
	return &apps.Manifest{
		Common: apps.Common{
			AppID:       AppID,
			Type:        apps.AppTypeHTTP,
			DisplayName: AppDisplayName,
			Description: AppDescription,
			HomepageURL: appURL(conf, "/"),
		},
		RootURL: appURL(conf, ""),
		RequestedPermissions: apps.Permissions{
			apps.PermissionUserJoinedChannelNotification,
			apps.PermissionActAsUser,
			apps.PermissionActAsBot,
		},
		RequestedLocations: apps.Locations{
			apps.LocationChannelHeader,
			apps.LocationPostMenu,
			apps.LocationCommand,
			apps.LocationInPost,
		},
	}
}

func (h *helloapp) handleManifest(w http.ResponseWriter, req *http.Request) {
	httputils.WriteJSON(w, Manifest(h.Conf.Get()))
}

func (h *helloapp) Install(call *apps.Call) *apps.CallResponse {
	return h.HelloApp.Install(AppID, AppDisplayName, call)
}

func handle(r *mux.Router, path string, h func(*apps.Call) *apps.CallResponse) {
	r.HandleFunc(path,
		func(w http.ResponseWriter, req *http.Request) {
			_, err := checkJWT(req)
			if err != nil {
				proxy.WriteCallError(w, http.StatusUnauthorized, err)
				return
			}

			call, err := apps.UnmarshalCallFromReader(req.Body)
			if err != nil {
				proxy.WriteCallError(w, http.StatusInternalServerError, err)
				return
			}

			httputils.WriteJSON(w, h(call))
		},
	).Methods("POST")
}

func checkJWT(req *http.Request) (*apps.JWTClaims, error) {
	authValue := req.Header.Get(apps.OutgoingAuthHeader)
	if !strings.HasPrefix(authValue, "Bearer ") {
		return nil, errors.Errorf("missing %s: Bearer header", apps.OutgoingAuthHeader)
	}

	jwtoken := strings.TrimPrefix(authValue, "Bearer ")
	claims := apps.JWTClaims{}
	_, err := jwt.ParseWithClaims(jwtoken, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(AppSecret), nil
	})
	if err != nil {
		return nil, err
	}

	return &claims, nil
}

func appURL(conf config.Config, path string) string {
	return conf.PluginURL + config.HelloHTTPPath + path
}
