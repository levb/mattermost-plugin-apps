package builtin

import (
	"github.com/mattermost/mattermost-plugin-apps/server/api"
)

func (a *App) funcGetBindings(call *api.Call) *api.CallResponse {
	return api.NewCallResponse("", a.bindings(call), nil)
}

func (a *App) bindings(call *api.Call) []*api.Binding {
	commands := []*api.Binding{
		{
			Label:       CommandInfo,
			Location:    CommandInfo,
			Description: "displays Apps plugin info",
			Call: &api.Call{
				URL: PathInfo,
			},
		}, {
			Label:       CommandConnect,
			Location:    CommandConnect,
			Hint:        "[AppID]",
			Description: "Connect an App to your Mattermost account",
			Call: &api.Call{
				URL: PathConnect,
			},
		}, {
			Label:       CommandDisconnect,
			Location:    CommandDisconnect,
			Hint:        "[AppID]",
			Description: "Disconnect an App from your Mattermost account",
			Call: &api.Call{
				URL: PathDisconnect,
			},
		},
	}

	var adminCommands []*api.Binding
	user, _ := a.API.Mattermost.User.Get(call.Context.ActingUserID)
	if user != nil && user.IsSystemAdmin() {
		adminCommands = []*api.Binding{
			{
				Label:       CommandDebug,
				Location:    CommandDebug,
				Hint:        "clean | install | view",
				Description: "debugging commands",
				Bindings: []*api.Binding{
					{
						Label:       CommandClean,
						Location:    CommandClean,
						Description: "clean the Apps KV store and config",
						Call: &api.Call{
							URL: PathDebugClean,
						},
					}, {
						Label:       CommandInstall,
						Location:    CommandInstall,
						Description: "install apps",
						Call: &api.Call{
							URL: PathDebugInstall,
						},
					}, {
						Label:       CommandBindings,
						Location:    CommandBindings,
						Description: "view bindings",
						Call: &api.Call{
							URL: PathDebugBindings,
						},
					},
				},
			},
		}
	}

	bindings := []*api.Binding{
		{
			Location: api.LocationCommand,
			Bindings: append(commands, adminCommands...),
		},
	}
	return bindings
}
