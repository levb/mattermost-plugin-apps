// Copyright (c) 2020-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package upawslambda

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/service/lambda"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/server/aws"
	"github.com/mattermost/mattermost-plugin-apps/server/utils"
)

type Upstream struct {
	app    *apps.App
	m      *apps.Manifest
	aws    aws.Service
	bucket string
}

func NewUpstream(app *apps.App, m *apps.Manifest, aws aws.Service, bucket string) *Upstream {
	return &Upstream{
		app:    app,
		aws:    aws,
		bucket: bucket,
	}
}

func (up *Upstream) OneWay(call *apps.Call) error {
	name := match(call.URL, up.m.FunctionRoutes)
	if name == "" {
		return utils.ErrNotFound
	}

	_, err := up.aws.Client().InvokeLambda(name, lambda.InvocationTypeEvent, call)
	return err
}

func (up *Upstream) Roundtrip(call *apps.Call) (io.ReadCloser, error) {
	name := match(call.URL, up.m.FunctionRoutes)
	if name == "" {
		return nil, utils.ErrNotFound
	}

	bb, err := up.aws.Client().InvokeLambda(name, lambda.InvocationTypeRequestResponse, call)
	if err != nil {
		return nil, err
	}
	return ioutil.NopCloser(bytes.NewReader(bb)), err
}

func (up *Upstream) GetStatic(path string) ([]byte, error) {
	name := match(path, up.m.StaticRoutes)
	if name == "" {
		return nil, utils.ErrNotFound
	}

	data, err := up.aws.Client().GetS3(up.bucket, name)
	if err != nil {
		return nil, err
	}
	return data, err
}

func match(callPath string, routes map[string]string) string {
	matchedName := ""
	matchedPath := ""
	for path, name := range routes {
		if strings.HasPrefix(callPath, path) {
			if len(path) > len(matchedPath) {
				matchedPath = path
				matchedName = name
			}
		}
	}
	return matchedName
}
