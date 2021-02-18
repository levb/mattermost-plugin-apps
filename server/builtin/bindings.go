package builtin

import (
	"github.com/mattermost/mattermost-plugin-apps/apps"
)

func (a *App) funcGetBindings(call *apps.Call) *apps.CallResponse {
	return apps.NewCallResponse("", a.bindings(call), nil)
}

func (a *App) bindings(call *apps.Call) []*apps.Binding {
	simple := func(label, path, hint, descr string) *apps.Binding {
		return &apps.Binding{
			Label:       label,
			Location:    apps.Location(label),
			Hint:        hint,
			Description: descr,
			Call: &apps.Call{
				URL: path,
			},
		}
	}

	commands := []*apps.Binding{
		simple(CommandInfo, PathInfo, "", "displays Apps plugin info"),
		simple(CommandList, PathList, "", "displays Apps plugin info"),
		simple(CommandConnect, PathConnect, "[AppID]", "Connect an App to your Mattermost account"),
		simple(CommandDisconnect, PathDisconnect, "[AppID]", "Disconnect an App from your Mattermost account"),
	}

	adminCommands := []*apps.Binding{
		simple(CommandInstall, PathInstallAppCommand, "[flags]", "Install an App to this Mattermost instance"),
		{
			Label:       CommandDebug,
			Location:    CommandDebug,
			Hint:        "clean | view",
			Description: "debugging commands",
			Bindings: []*apps.Binding{
				simple(CommandClean, PathDebugClean, "", "clean the Apps KV store and config"),
				// simple(CommandInstall, PathDebugInstall, "", "install apps"),
				simple(CommandBindings, PathDebugBindings, "", "view bindings"),
			},
		},
	}

	user, _ := a.mm.User.Get(call.Context.ActingUserID)
	if user != nil && user.IsSystemAdmin() {
		commands = append(commands, adminCommands...)
	}

	return []*apps.Binding{
		{
			Location: apps.LocationCommand,
			Bindings: commands,
		},
	}
}
