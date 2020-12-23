package builtin

import (
	"github.com/mattermost/mattermost-plugin-apps/server/api"
	"github.com/mattermost/mattermost-plugin-apps/server/utils/md"
)

func (a *App) connectForm(c *api.Call) (*api.Form, error) {
	apps, err := a.API.Admin.ListApps()
	if err != nil {
		return nil, err
	}
	return &api.Form{
		Title: "Connect App to Mattermost",
		Fields: []*api.Field{
			{
				Name:                 fieldAppID,
				Type:                 api.FieldTypeStaticSelect,
				Description:          "App to connect to",
				Label:                fieldAppID,
				AutocompleteHint:     "enter or select an App to connect",
				AutocompletePosition: 1,
				SelectStaticOptions:  a.getAppOptions(apps),
			},
		},
	}, nil
}

func (a *App) connect(call *api.Call) (md.MD, error) {
	appID := api.AppID(call.GetValue(fieldAppID, ""))

	connectURL, err := a.API.Proxy.StartOAuthConnect(call.Context.ActingUserID, appID, call)
	if err != nil {
		return "", err
	}

	return md.Markdownf("Click [here](%s) to continue", connectURL), nil
}

func (a *App) getAppOptions(apps []*api.App) []api.SelectOption {
	options := []api.SelectOption{}
	for _, app := range apps {
		if app.GrantedPermissions.Contains(api.PermissionActAsUser) {
			options = append(options, api.SelectOption{
				Label: app.Manifest.DisplayName,
				Value: string(app.Manifest.AppID),
			})
		}
	}
	return options
}
