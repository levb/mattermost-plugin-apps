// Copyright (c) 2020-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package uphttp

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/server/httpout"
	"github.com/mattermost/mattermost-plugin-apps/upstream"
	"github.com/mattermost/mattermost-plugin-apps/utils"
	"github.com/mattermost/mattermost-plugin-apps/utils/httputils"
)

type Upstream struct {
	StaticUpstream
	appSecret string
}

var _ upstream.Upstream = (*Upstream)(nil)

func NewUpstream(app *apps.App, httpOut httpout.Service) *Upstream {
	staticUp := NewStaticUpstream(&app.Manifest, httpOut)
	return &Upstream{
		StaticUpstream: *staticUp,
		appSecret:      app.Secret,
	}
}

func (u *Upstream) Roundtrip(call *apps.CallRequest, async bool) (io.ReadCloser, error) {
	if async {
		go func() {
			resp, _ := u.invoke(call.Context.BotUserID, call)
			if resp != nil {
				resp.Body.Close()
			}
		}()
		return nil, nil
	}

	resp, err := u.invoke(call.Context.ActingUserID, call) // nolint:bodyclose
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (u *Upstream) invoke(fromMattermostUserID string, call *apps.CallRequest) (*http.Response, error) {
	if call == nil {
		return nil, utils.NewInvalidError("empty call")
	}

	return u.post(call.Context.ActingUserID, u.rootURL+call.Path, call)
}

// post does not close resp.Body, it's the caller's responsibility
func (u *Upstream) post(fromMattermostUserID string, url string, msg interface{}) (*http.Response, error) {
	client := u.httpOut.MakeClient(true)
	jwtoken, err := createJWT(fromMattermostUserID, u.appSecret)
	if err != nil {
		return nil, err
	}

	piper, pipew := io.Pipe()
	go func() {
		encodeErr := json.NewEncoder(pipew).Encode(msg)
		if encodeErr != nil {
			_ = pipew.CloseWithError(encodeErr)
		}
		pipew.Close()
	}()

	req, err := http.NewRequest(http.MethodPost, url, piper)
	if err != nil {
		return nil, err
	}
	req.Header.Set(apps.OutgoingAuthHeader, "Bearer "+jwtoken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		bb, _ := httputils.ReadAndClose(resp.Body)
		return nil, errors.New(string(bb))
	}

	return resp, nil
}

func createJWT(actingUserID, secret string) (string, error) {
	claims := apps.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
		ActingUserID: actingUserID,
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}
