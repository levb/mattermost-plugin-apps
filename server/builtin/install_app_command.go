// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package builtin

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/server/examples/go/hello/builtin_hello"
	"github.com/mattermost/mattermost-plugin-apps/server/examples/go/hello/http_hello"
	"github.com/mattermost/mattermost-plugin-apps/server/examples/js/aws_hello"
)

const (
	DebugInstallFromURL = true
)

func (a *App) installAppCommandForm(c *apps.Call) (*apps.Form, error) {
	fields := []*apps.Field{
		{
			Name:             fieldAppID,
			Type:             apps.FieldTypeStaticSelect,
			Description:      "select an App from the list",
			Label:            flagAppID,
			AutocompleteHint: "App",
			SelectStaticOptions: []apps.SelectOption{
				{
					Label: builtin_hello.AppDisplayName,
					Value: builtin_hello.AppID,
				}, {
					Label: http_hello.AppDisplayName,
					Value: http_hello.AppID,
				}, {
					Label: aws_hello.AppDisplayName,
					Value: aws_hello.AppID,
				},
			},
		},
	}

	if DebugInstallFromURL {
		fields = append(fields, &apps.Field{
			Name:             fieldManifestURL,
			Type:             apps.FieldTypeText,
			Description:      "(debug) location of the App manifest",
			Label:            flagManifestURL,
			AutocompleteHint: "enter the URL",
		})
	}

	fmt.Printf("<><> %+v\n", fields[0])

	return &apps.Form{
		Title:  "Install an App",
		Fields: fields,
		Call: &apps.Call{
			URL: PathInstallAppCommand,
		},
	}, nil
}

func (a *App) installAppCommand(call *apps.Call) *apps.CallResponse {
	// id := call.GetStringValue(fieldAppID, "")
	// manifestURL := call.GetStringValue(fieldManifestURL, "")
	// conf := a.conf.Get()

	// var manifest *apps.Manifest
	// switch {
	// case id != "" && manifestURL != "":
	// 	return apps.NewCallResponse("", nil,
	// 		errors.Errorf("`--%s` and `--%s` can not be both specified", flagAppID, flagManifestURL))
	// case id == http_hello.AppID:
	// 	manifest = http_hello.Manifest(conf)
	// case id == builtin_hello.AppID:
	// 	manifest = builtin_hello.Manifest()
	// case id == aws_hello.AppID:
	// 	manifest = aws_hello.Manifest()
	// case manifestURL != "":
	// var err error
	// manifest, err = proxy.LoadManifest(manifestURL)
	// if err != nil {
	// 	return apps.NewCallResponse("", nil, err)
	// }
	// }

	return &apps.CallResponse{
		// Type: apps.CallResponseTypeForm,
		// Form: a.installAppForm(manifest),
	}
}
