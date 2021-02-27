package restapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-apps/apps"
)

func (a *restapi) handleSubscribe(w http.ResponseWriter, r *http.Request) {
	a.handleSubscribeCore(w, r, true)
}

func (a *restapi) handleUnsubscribe(w http.ResponseWriter, r *http.Request) {
	a.handleSubscribeCore(w, r, false)
}

func (a *restapi) handleSubscribeCore(w http.ResponseWriter, r *http.Request, isSubscribe bool) {
	var err error
	actingUserID := ""
	// logMessage := ""
	status := http.StatusOK

	defer func() {
		resp := apps.SubscriptionResponse{}
		if err != nil {
			resp.Error = errors.Wrap(err, "failed operation").Error()
			status = http.StatusInternalServerError
			// logMessage = "Error: " + resp.Error
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write(resp.ToJSON())
	}()

	fmt.Printf("<><> handleSubscribeCore: 1\n")
	actingUserID = r.Header.Get("Mattermost-User-ID")

	if actingUserID == "" {
		fmt.Printf("<><> handleSubscribeCore: 2 not logged in\n")
		err = errors.New("user not logged in")
		status = http.StatusUnauthorized
		return
	}

	// TODO check for sysadmin
	if !a.mm.User.HasPermissionTo(actingUserID, model.PERMISSION_MANAGE_SYSTEM) {
		fmt.Printf("<><> handleSubscribeCore: 3 not sysadmin\n")
		http.Error(w, errors.New("forbidden").Error(), http.StatusForbidden)
		return
	}

	var sub apps.Subscription
	if err = json.NewDecoder(r.Body).Decode(&sub); err != nil {
		fmt.Printf("<><> handleSubscribeCore: 4 failed to decode\n")
		status = http.StatusBadRequest
		return
	}

	// TODO replace with an appropriate API-level call that would validate,
	// deduplicate, etc.
	if isSubscribe {
		err = a.appservices.Subscribe(&sub)
	} else {
		err = a.appservices.Unsubscribe(&sub)
	}
	fmt.Printf("<><> handleSubscribeCore: 5 %v\n", err)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
