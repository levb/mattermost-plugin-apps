package builtin

import (
	"github.com/mattermost/mattermost-plugin-apps/server/api"
)

func (a *App) funcGetBindings(call *api.Call) *api.CallResponse {
	return api.NewCallResponse("", a.bindings(call), nil)
}

func (a *App) bindings(call *api.Call) []*api.Binding {
	simple := func(label, path, hint, descr string) *api.Binding {
		return &api.Binding{
			Label:       label,
			Location:    api.Location(label),
			Hint:        hint,
			Description: descr,
			Call: &api.Call{
				URL: path,
			},
		}
	}

	commands := []*api.Binding{
		simple(CommandInfo, PathInfo, "", "displays Apps plugin info"),
		simple(CommandList, PathList, "", "displays Apps plugin info"),
		simple(CommandConnect, PathConnect, "[AppID]", "Connect an App to your Mattermost account"),
		simple(CommandDisconnect, PathDisconnect, "[AppID]", "Disconnect an App from your Mattermost account"),
	}

	adminCommands := []*api.Binding{
		simple(CommandInstall, PathInstallAppCommand, "[flags]", "Install an App to this Mattermost instance"),
		{
			Label:       CommandDebug,
			Location:    CommandDebug,
			Hint:        "clean | view",
			Description: "debugging commands",
			Bindings: []*api.Binding{
				simple(CommandClean, PathDebugClean, "", "clean the Apps KV store and config"),
				// simple(CommandInstall, PathDebugInstall, "", "install apps"),
				simple(CommandBindings, PathDebugBindings, "", "view bindings"),
			},
		},
	}

	user, _ := a.API.Mattermost.User.Get(call.Context.ActingUserID)
	if user != nil && user.IsSystemAdmin() {
		commands = append(commands, adminCommands...)
	}

	return []*api.Binding{
		{
			Location: api.LocationCommand,
			Bindings: commands,
		},
	}
}
