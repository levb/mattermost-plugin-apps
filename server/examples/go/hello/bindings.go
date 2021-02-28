package hello

import (
	"github.com/mattermost/mattermost-plugin-apps/apps"
)

func (h *HelloApp) GetBindings(call *apps.Call) *apps.CallResponse {
	return apps.NewCallResponse("", bindings(call), nil)
}

func bindings(call *apps.Call) []*apps.Binding {
	post := apps.MakeCall(PathPostAsUser)
	post.Expand = &apps.Expand{
		ActingUserAccessToken: apps.ExpandAll,
	}

	commandToModal := apps.MakeCall(PathSendSurveyCommandToModal)

	justSend := apps.MakeCall(PathSendSurvey)

	modalFromPost := apps.MakeCall(PathSendSurveyModal)
	modalFromPost.Expand = &apps.Expand{Post: apps.ExpandAll}

	channelHeaderBindings := []*apps.Binding{
		{
			Location:    "send",
			Label:       "Survey a user",
			Icon:        "https://raw.githubusercontent.com/mattermost/mattermost-plugin-jira/master/assets/icon.svg",
			Hint:        "Send survey to a user",
			Description: "Send a customized emotional response survey to a user",
			Call: &apps.Call{
				URL: PathSendSurvey,
			},
		},
	}

	postMenuBindings := []*apps.Binding{
		{
			Location:    "send-me",
			Label:       "Survey myself",
			Hint:        "Send survey to myself",
			Description: "Send a customized emotional response survey to myself",
			Call:        justSend, // will use ActingUserID by default
		},
		{
			Location:    "send",
			Label:       "Survey a user",
			Hint:        "Send survey to a user",
			Description: "Send a customized emotional response survey to a user",
			Call:        modalFromPost,
		},
	}

	// TODO /Command binding is a placeholder, may not be final, test!
	commandBindings := []*apps.Binding{
		{
			Label:       "message",
			Location:    "message",
			Hint:        "[@user] [--message text]",
			Description: "send a message to a user",
			Call:        justSend,
		}, {
			Label:       "message-modal",
			Location:    "message-modal",
			Hint:        "[--message] message",
			Description: "send a message to a user",
			Call:        commandToModal,
		}, {
			Label:       "manage",
			Location:    "manage",
			Hint:        "subscribe | unsubscribe ",
			Description: "manage channel subscriptions to greet new users",
			Bindings: []*apps.Binding{
				{
					Label:       "subscribe",
					Location:    "subscribe",
					Hint:        "[--channel]",
					Description: "subscribes a channel to greet new users",
					Call:        apps.MakeCall(PathSubscribeChannel, "mode", "on"),
				}, {
					Label:       "unsubscribe",
					Location:    "unsubscribe",
					Hint:        "[--channel]",
					Description: "unsubscribes a channel from greeting new users",
					Call:        apps.MakeCall(PathSubscribeChannel, "mode", "off"),
				},
			},
		},
	}

	// TODO /Command binding is a placeholder, may not be final, test!
	connectedCommandBindings := []*apps.Binding{}
	// if call.Context.ActingUserIsConnected {
	connectedCommandBindings = []*apps.Binding{
		{
			Label:         "post",
			Location:      "post",
			Hint:          "--message text",
			Description:   "post a message to channel user",
			DependsOnUser: true,
			Call:          post,
		},
	}
	// }

	return []*apps.Binding{
		{
			Location: apps.LocationCommand,
			Bindings: append(commandBindings, connectedCommandBindings...),
		},
		{
			// TODO make this a subscribe button, with a state (current subscription status)
			Location: apps.LocationChannelHeader,
			Bindings: channelHeaderBindings,
		},
		{
			Location: apps.LocationPostMenu,
			Bindings: postMenuBindings,
		},
	}
}
