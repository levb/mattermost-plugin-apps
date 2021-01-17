package http_hello

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-apps/server/api"
	"github.com/mattermost/mattermost-plugin-apps/server/api/impl/proxy"
	"github.com/mattermost/mattermost-plugin-apps/server/examples/go/hello"
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

// Init hello app router
func Init(router *mux.Router, appsService *api.Service) {
	h := helloapp{
		hello.NewHelloApp(appsService),
	}

	r := router.PathPrefix(api.HelloHTTPPath).Subrouter()
	r.HandleFunc(PathManifest, h.handleManifest).Methods("GET")

	handle(r, api.DefaultInstallCallPath, h.Install)
	handle(r, api.DefaultBindingsCallPath, h.GetBindings)
	handle(r, hello.PathSendSurvey, h.SendSurvey)
	handle(r, hello.PathSendSurveyModal, h.SendSurveyModal)
	handle(r, hello.PathSendSurveyCommandToModal, h.SendSurveyCommandToModal)
	handle(r, hello.PathSurvey, h.Survey)
	handle(r, hello.PathUserJoinedChannel, h.UserJoinedChannel)
	handle(r, hello.PathPostAsUser, h.PostAsUser)
}

func (h *helloapp) handleManifest(w http.ResponseWriter, req *http.Request) {
	httputils.WriteJSON(w,
		api.Manifest{
			AppID:       AppID,
			Type:        api.AppTypeHTTP,
			DisplayName: AppDisplayName,
			Description: AppDescription,
			HTTPRootURL: h.appURL(""),
			RequestedPermissions: api.Permissions{
				api.PermissionUserJoinedChannelNotification,
				api.PermissionActAsUser,
				api.PermissionActAsBot,
			},
			RequestedLocations: api.Locations{
				api.LocationChannelHeader,
				api.LocationPostMenu,
				api.LocationCommand,
				api.LocationInPost,
			},
			HomepageURL: h.appURL("/"),
		})
}

func (h *helloapp) Install(call *api.Call) *api.CallResponse {
	return h.HelloApp.Install(AppID, AppDisplayName, call)
}

func handle(r *mux.Router, path string, h func(*api.Call) *api.CallResponse) {
	r.HandleFunc(path,
		func(w http.ResponseWriter, req *http.Request) {
			_, err := checkJWT(req)
			if err != nil {
				proxy.WriteCallError(w, http.StatusUnauthorized, err)
				return
			}

			call, err := api.UnmarshalCallFromReader(req.Body)
			if err != nil {
				proxy.WriteCallError(w, http.StatusInternalServerError, err)
				return
			}

			httputils.WriteJSON(w, h(call))
		},
	).Methods("POST")
}

func checkJWT(req *http.Request) (*api.JWTClaims, error) {
	authValue := req.Header.Get(api.OutgoingAuthHeader)
	if !strings.HasPrefix(authValue, "Bearer ") {
		return nil, errors.Errorf("missing %s: Bearer header", api.OutgoingAuthHeader)
	}

	jwtoken := strings.TrimPrefix(authValue, "Bearer ")
	claims := api.JWTClaims{}
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

func (h *helloapp) appURL(path string) string {
	conf := h.API.Configurator.GetConfig()
	return conf.PluginURL + api.HelloHTTPPath + path
}
