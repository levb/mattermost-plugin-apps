// Copyright (c) 2020-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package upawslambda

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/service/lambda"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/server/aws"
	"github.com/mattermost/mattermost-plugin-apps/server/utils"
)

type Upstream struct {
	app *apps.App
	m   *apps.Manifest
	aws aws.Service
}

func NewUpstream(app *apps.App, m *apps.Manifest, aws aws.Service) *Upstream {
	return &Upstream{
		app: app,
		aws: aws,
	}
}

func (up *Upstream) OneWay(call *apps.Call) error {
	name, ok := up.m.LambdaRoutes[call.URL]
	if !ok {
		return utils.ErrNotFound
	}

	_, err := up.aws.Client().InvokeLambda(name, lambda.InvocationTypeEvent, call)
	return err
}

func (up *Upstream) Roundtrip(call *apps.Call) (io.ReadCloser, error) {
	name, ok := up.m.LambdaRoutes[call.URL]
	if !ok {
		return nil, utils.ErrNotFound
	}

	bb, err := up.aws.Client().InvokeLambda(name, lambda.InvocationTypeRequestResponse, call)
	if err != nil {
		return nil, err
	}
	return ioutil.NopCloser(bytes.NewReader(bb)), err
}
