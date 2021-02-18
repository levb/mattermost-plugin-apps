package builtin

import (
	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/server/utils/md"
)

func (a *App) connectForm(c *apps.Call) (*apps.Form, error) {
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
				SelectStaticOptions:  a.getAppOptions(),
			},
		},
	}, nil
}

func (a *App) connect(call *apps.Call) *apps.CallResponse {
	appID := apps.AppID(call.GetStringValue(fieldAppID, ""))
	connectURL, err := a.proxy.StartOAuthConnect(call.Context.ActingUserID, appID, call)
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

func (a *App) getAppOptions() []apps.SelectOption {
	options := []apps.SelectOption{}

	allApps := a.proxy.ListInstalledApps()
	for _, app := range allApps {
		if app.GrantedPermissions.Contains(apps.PermissionActAsUser) {
			options = append(options, apps.SelectOption{
				Label: app.DisplayName,
				Value: string(app.AppID),
			})
		}
	}
	return options
}
