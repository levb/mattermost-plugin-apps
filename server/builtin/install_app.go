// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package builtin

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/server/store"
	"github.com/mattermost/mattermost-plugin-apps/server/utils/httputils"
	"github.com/pkg/errors"
)

const (
	DebugInstallFromURL = true
)

func (a *builtinApp) installMarketplaceCommandForm(c *apps.Call) *apps.CallResponse {
	return &apps.CallResponse{
		Type: apps.CallResponseTypeForm,
		Form: &apps.Form{
			Title: "Install an App from marketplace",
			Fields: []*apps.Field{
				{
					Name:                 fAppID,
					Type:                 apps.FieldTypeDynamicSelect,
					Description:          "select a Marketplace App",
					Label:                fAppID,
					AutocompleteHint:     "App ID",
					AutocompletePosition: 1,
				},
			},
			Call: &apps.Call{
				URL: PathInstallMarketplace,
			},
		},
	}
}

func (a *builtinApp) installDeveloperCommandForm(c *apps.Call) *apps.CallResponse {
	return &apps.CallResponse{
		Type: apps.CallResponseTypeForm,
		Form: &apps.Form{
			Title: "Install an App",
			Fields: []*apps.Field{
				{
					Name:                 fURL,
					Type:                 apps.FieldTypeText,
					Description:          "URL of the App manifest",
					Label:                fURL,
					AutocompleteHint:     "enter the URL",
					AutocompletePosition: 1,
				},
			},
			Call: &apps.Call{
				URL: PathInstallDeveloper,
			},
		},
	}
}

func (a *builtinApp) installMarketplaceCommand(call *apps.Call) *apps.CallResponse {
	appID := apps.AppID(call.GetStringValue(fAppID, ""))
	m, err := a.store.Manifest.Get(appID)
	if err != nil {
		return apps.NewErrorCallResponse(err)
	}

	return a.installAppFormManifest(m, call)
}

func (a *builtinApp) installDeveloperCommand(call *apps.Call) *apps.CallResponse {
	url := call.GetStringValue(fURL, "")
	data, err := httputils.GetFromURL(url)
	if err != nil {
		return apps.NewErrorCallResponse(err)
	}
	m, err := store.DecodeManifest(data)
	if err != nil {
		return apps.NewErrorCallResponse(err)
	}
	err = a.store.Manifest.StoreLocal(m)
	if err != nil {
		return apps.NewErrorCallResponse(err)
	}

	return a.installAppFormManifest(m, call)
}

func (a *builtinApp) installMarketplaceLookup(call *apps.Call) []*apps.SelectOption {
	name := call.GetStringValue("name", "")
	input := call.GetStringValue("user_input", "")

	switch name {
	case fAppID:
		marketplaceApps := a.proxy.ListMarketplaceApps(input)
		var options []*apps.SelectOption
		for _, mapp := range marketplaceApps {
			if !mapp.Installed {
				options = append(options, &apps.SelectOption{
					Value: string(mapp.Manifest.AppID),
					Label: mapp.Manifest.DisplayName,
				})
			}
		}
		return options
	}
	return nil
}

func (a *builtinApp) installAppForm(c *apps.Call) *apps.CallResponse {
	s, ok := c.Context.Props[contextInstallAppID]
	if !ok {
		return apps.NewErrorCallResponse(errors.New("no AppID to install in Context"))
	}
	appID := apps.AppID(s)

	m, err := a.store.Manifest.Get(appID)
	if err != nil {
		return apps.NewErrorCallResponse(err)
	}

	return a.installAppFormManifest(m, c)
}

func (a *builtinApp) installAppFormManifest(m *apps.Manifest, c *apps.Call) *apps.CallResponse {
	fields := []*apps.Field{}

	if len(m.RequestedLocations) > 0 {
		fields = append(fields, &apps.Field{
			Name:        fConsentLocations,
			Type:        apps.FieldTypeBool,
			Label:       "Application may display its UI elements in the following locations",
			Description: fmt.Sprintf("%s", m.RequestedLocations),
			IsRequired:  true,
		})
	}

	if len(m.RequestedPermissions) > 0 {
		fields = append(fields, &apps.Field{
			Name:        fConsentPermissions,
			Type:        apps.FieldTypeBool,
			Label:       "Application will have the following permissions",
			Description: fmt.Sprintf("%s", m.RequestedPermissions),
			IsRequired:  true,
		})
	}

	fields = append(fields, &apps.Field{
		Name:        fRequireUserConsent,
		Type:        apps.FieldTypeBool,
		Label:       fmt.Sprintf("Require explicit user's consent to allow %s App impersonate the user", m.AppID),
		Description: "If off, users will be quietly connected to the App as needed; otherwise prompt for consent.",
	})

	if m.Type == apps.AppTypeHTTP {
		fields = append(fields, &apps.Field{
			Name:             fSecret,
			Type:             apps.FieldTypeText,
			Description:      "The App's secret to use in JWT.",
			Label:            fSecret,
			AutocompleteHint: "paste the secret obtained from the App",
			IsRequired:       true,
		})
	}

	cr := &apps.CallResponse{
		Type: apps.CallResponseTypeForm,
		Form: &apps.Form{
			Title:  fmt.Sprintf("Install App %s", m.DisplayName),
			Fields: fields,
			Call: &apps.Call{
				URL: PathInstallApp,
				Context: &apps.Context{
					Props: map[string]string{
						contextInstallAppID: string(m.AppID),
					},
				},
				Values: map[string]interface{}{
					fAppID: string(m.AppID),
				},
				Expand: &apps.Expand{
					AdminAccessToken: apps.ExpandAll,
				},
			},
		},
	}
	return cr
}

func (a *builtinApp) installApp(call *apps.Call) *apps.CallResponse {
	secret := call.GetStringValue(fSecret, "")
	requireUserConsent := call.GetBoolValue(fRequireUserConsent)
	locationsConsent := call.GetBoolValue(fConsentLocations)
	permissionsConsent := call.GetBoolValue(fConsentPermissions)
	id := ""
	if v, _ := call.Context.Props[contextInstallAppID]; v != "" {
		id = v
	}

	m, err := a.store.Manifest.Get(apps.AppID(id))
	if err != nil {
		return apps.NewErrorCallResponse(errors.Wrap(err, "failed to load App manifest"))
	}

	if !locationsConsent && len(m.RequestedLocations) > 0 {
		return apps.NewErrorCallResponse(errors.New("consent to grant access to UI locations is required to install"))
	}
	if !permissionsConsent && len(m.RequestedPermissions) > 0 {
		return apps.NewErrorCallResponse(errors.New("consent to grant permissions is required to install"))
	}

	app, out, err := a.proxy.InstallApp(call.Context, &apps.InInstallApp{
		Manifest:         m,
		OAuth2TrustedApp: !requireUserConsent,
		Secret:           secret,
	})
	return apps.NewCallResponse(out, app, err)
}
