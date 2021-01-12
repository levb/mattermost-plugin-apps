// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package builtin

import (
	"github.com/mattermost/mattermost-plugin-apps/server/api"
	"github.com/mattermost/mattermost-plugin-apps/server/utils/md"
)

func (a *App) installCommandForm(c *api.Call) (*api.Form, error) {
	f := &api.Form{
		Title: "Set up an App to Mattermost",
		Fields: []*api.Field{
			{
				Name:             fieldManifestURL,
				Type:             api.FieldTypeText,
				Description:      "(debug) location of the App manifest",
				Label:            flagManifestURL,
				AutocompleteHint: "enter the URL",
				IsRequired:       true,
			},
		},
	}

			{
				Name:             fieldSecret,
				Type:             api.FieldTypeText,
				Description:      "The App's secret to use in JWT.",
				Label:            flagSecret,
				AutocompleteHint: "paste the secret obtained from the App",
			},
			{
				Name:        fieldRequireUserConsent,
				Type:        api.FieldTypeBool,
				Description: "If **on**, users will be prompted for consent before connecting to the App",
				Label:       flagRequireUserConsent,
			},
		},
	}, nil
}

func (a *App) install(call *api.Call) (md.MD, error) {
	appTy := call.GetValue(fieldManifestURL, "")
	manifestURL := call.GetValue(fieldManifestURL, "")
	secret := call.GetValue(fieldSecret, "")
	requireUserConsent := call.GetValueBool(fieldRequireUserConsent)

	if 







	return md.Markdownf("Click [here](%s) to continue", connectURL), nil
}
