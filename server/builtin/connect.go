package builtin

import (
	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/server/utils/md"
)

func (a *App) connectForm(c *apps.Call) (*apps.Form, error) {
	all, err := a.API.Admin.ListApps()
	if err != nil {
		return nil, err
	}
	return &apps.Form{
		Title: "Connect App to Mattermost",
		Fields: []*apps.Field{
			{
				Name:                 fieldAppID,
				Type:                 apps.FieldTypeStaticSelect,
				Description:          "App to connect to",
				Label:                flagAppID,
				AutocompleteHint:     "enter or select an App to connect",
				AutocompletePosition: 1,
				SelectStaticOptions:  a.getAppOptions(all),
			},
		},
	}, nil
}

func (a *App) connect(call *apps.Call) *apps.CallResponse {
	appID := apps.AppID(call.GetStringValue(fieldAppID, ""))
	connectURL, err := a.API.Proxy.StartOAuthConnect(call.Context.ActingUserID, appID, call)
	if err != nil {
		return apps.NewCallResponse("", nil, err)
	}
	txt := md.Markdownf("Click [here](%s) to continue", connectURL)
	return apps.NewCallResponse(txt, nil, nil)
}

func (a *App) disconnect(call *apps.Call) *apps.CallResponse {
	appID := apps.AppID(call.GetStringValue(fieldAppID, ""))
	txt := md.Markdownf("TODO: disconnect %s", appID)
	return apps.NewCallResponse(txt, nil, nil)
}

func (a *App) getAppOptions(all []*apps.App) []apps.SelectOption {
	options := []apps.SelectOption{}
	for _, app := range all {
		if app.GrantedPermissions.Contains(apps.PermissionActAsUser) {
			options = append(options, apps.SelectOption{
				Label: app.Manifest.DisplayName,
				Value: string(app.Manifest.AppID),
			})
		}
	}
	return options
}
