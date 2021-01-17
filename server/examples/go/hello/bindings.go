package hello

import (
	"github.com/mattermost/mattermost-plugin-apps/server/api"
)

func (h *HelloApp) GetBindings(call *api.Call) *api.CallResponse {
	return api.NewCallResponse("", bindings(call), nil)
}

func bindings(call *api.Call) []*api.Binding {
	post := api.MakeCall(PathPostAsUser)
	post.Expand = &api.Expand{
		ActingUserAccessToken: api.ExpandAll,
	}

	commandToModal := api.MakeCall(PathSendSurveyCommandToModal)

	justSend := api.MakeCall(PathSendSurvey)

	modalFromPost := api.MakeCall(PathSendSurveyModal)
	modalFromPost.Expand = &api.Expand{Post: api.ExpandAll}

	channelHeaderBindings := []*api.Binding{
		{
			Location:    "send",
			Label:       "Survey a user",
			Icon:        "https://raw.githubusercontent.com/mattermost/mattermost-plugin-jira/master/assets/icon.svg",
			Hint:        "Send survey to a user",
			Description: "Send a customized emotional response survey to a user",
			Call: &api.Call{
				URL: PathSendSurvey,
			},
		},
	}

	postMenuBindings := []*api.Binding{
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
	commandBindings := []*api.Binding{
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
			Bindings: []*api.Binding{
				{
					Label:       "subscribe",
					Location:    "subscribe",
					Hint:        "[--channel]",
					Description: "subscribes a channel to greet new users",
					Call:        api.MakeCall(PathSubscribeChannel, "mode", "on"),
				}, {
					Label:       "unsubscribe",
					Location:    "unsubscribe",
					Hint:        "[--channel]",
					Description: "unsubscribes a channel from greeting new users",
					Call:        api.MakeCall(PathSubscribeChannel, "mode", "off"),
				},
			},
		},
	}

	// TODO /Command binding is a placeholder, may not be final, test!
	connectedCommandBindings := []*api.Binding{}
	// if call.Context.ActingUserIsConnected {
	connectedCommandBindings = []*api.Binding{
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

	return []*api.Binding{
		{
			Location: api.LocationCommand,
			Bindings: append(commandBindings, connectedCommandBindings...),
		},
		{
			// TODO make this a subscribe button, with a state (current subscription status)
			Location: api.LocationChannelHeader,
			Bindings: channelHeaderBindings,
		},
		{
			Location: api.LocationPostMenu,
			Bindings: postMenuBindings,
		},
	}
}
