// Copyright (c) 2020-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package upaws

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/upstream"
	"github.com/mattermost/mattermost-plugin-apps/utils"
)

// Upstream wraps an awsClient to make requests to the App. It should not be
// reused between requests, nor cached.
type Upstream struct {
	StaticUpstream
	app *apps.App
}

var _ upstream.Upstream = (*Upstream)(nil)

// invocationPayload is a scoped down version of
// https://pkg.go.dev/github.com/aws/aws-lambda-go@v1.13.3/events#APIGatewayProxyRequest
type invocationPayload struct {
	Path       string            `json:"path"`
	HTTPMethod string            `json:"httpMethod"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

// invocationResponse is a scoped down version of
// https://pkg.go.dev/github.com/aws/aws-lambda-go@v1.13.3/events#APIGatewayProxyResponse
type invocationResponse struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

func NewUpstream(app *apps.App, awsClient Client, bucket string) *Upstream {
	staticUp := NewStaticUpstream(&app.Manifest, awsClient, bucket)
	return &Upstream{
		StaticUpstream: *staticUp,
	}
}

func (u *Upstream) Roundtrip(call *apps.CallRequest, async bool) (io.ReadCloser, error) {
	name := match(call.Path, u.manifest)
	if name == "" {
		return nil, utils.ErrNotFound
	}

	crString, err := u.InvokeFunction(name, async, call)
	if err != nil {
		return nil, err
	}
	return io.NopCloser(strings.NewReader(crString)), nil
}

// InvokeFunction is a public method used in appsctl, but is not a part of the
// upstream.Upstream interface. It invokes a function with a specified name,
// with no conversion.
func (u *Upstream) InvokeFunction(name string, async bool, call *apps.CallRequest) (string, error) {
	typ := lambda.InvocationTypeRequestResponse
	if async {
		typ = lambda.InvocationTypeEvent
	}

	payload, err := callToInvocationPayload(call)
	if err != nil {
		return "", errors.Wrap(err, "failed to convert call into invocation payload")
	}

	bb, err := u.awsClient.InvokeLambda(name, typ, payload)
	if async || err != nil {
		return "", err
	}

	var resp invocationResponse
	err = json.Unmarshal(bb, &resp)
	if err != nil {
		return "", errors.Wrap(err, "Error marshaling request payload")
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("lambda invocation failed with status code %v and body %v", resp.StatusCode, resp.Body)
	}

	return resp.Body, nil
}

func (u *Upstream) GetStatic(path string) (io.ReadCloser, int, error) {
	key := S3StaticName(u.app.AppID, u.app.Version, path)
	data, err := u.awsClient.GetS3(u.bucket, key)
	if err != nil {
		return nil, http.StatusBadRequest, errors.Wrapf(err, "can't download from S3:bucket:%s, path:%s", u.bucket, path)
	}
	return io.NopCloser(bytes.NewReader(data)), http.StatusOK, nil
}

func callToInvocationPayload(call *apps.CallRequest) ([]byte, error) {
	body, err := json.Marshal(call)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal call for lambda payload")
	}

	request := invocationPayload{
		Path:       call.Path,
		HTTPMethod: http.MethodPost,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(body),
	}

	payload, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal lambda payload")
	}

	return payload, nil
}

func match(callPath string, m *apps.Manifest) string {
	matchedName := ""
	matchedPath := ""
	for _, f := range m.AWSLambda {
		if strings.HasPrefix(callPath, f.Path) {
			if len(f.Path) > len(matchedPath) {
				matchedName = LambdaName(m.AppID, m.Version, f.Name)
				matchedPath = f.Path
			}
		}
	}

	return matchedName
}
