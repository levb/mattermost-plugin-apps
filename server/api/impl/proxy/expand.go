package proxy

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-apps/server/api"
	"github.com/mattermost/mattermost-server/v5/model"
)

type expander struct {
	// Context to expand (can be expanded multiple times on the same expander,
	// for different Apps).
	*api.Context

	sessionToken api.SessionToken
}

func (p *Proxy) newExpander(cc *api.Context, sessionToken api.SessionToken) *expander {
	e := &expander{
		Context:      cc,
		sessionToken: sessionToken,
	}
	return e
}

func (p *Proxy) ExpandForApp(e *expander, call *api.Call, app *api.App) (*api.Context, string, error) {
	expand := call.Expand
	cc := *e.Context
	fmt.Printf("<><> ExpandForApp 0: %v\n", call.URL)
	fmt.Printf("<><> ExpandForApp 1: %+v\n", expand)
	fmt.Printf("<><> ExpandForApp 2: %+v\n", cc)
	if expand == nil {
		expand = &api.Expand{}
	}

	cc.App = nil
	cc.BotUserID = app.BotUserID

	if e.MattermostSiteURL == "" {
		mmconf := p.conf.GetMattermostConfig()
		if mmconf.ServiceSettings.SiteURL != nil {
			e.MattermostSiteURL = *mmconf.ServiceSettings.SiteURL
		}
	}
	cc.MattermostSiteURL = e.MattermostSiteURL

	if app.GrantedPermissions.Contains(api.PermissionActAsBot) {
		cc.BotAccessToken = app.BotAccessToken
	}

	// TODO: do we need another way of obtaining an admin token? Should it be
	// out of expand?
	// TODO: Implement collecting user consent for the admin token, in-line?
	if expand.AdminAccessToken.Any() {
		cc.AdminAccessToken = string(e.sessionToken)
	}

	if e.ActingUserID != "" {
		if app.GrantedPermissions.Contains(api.PermissionActAsUser) {
			oauth := p.newMattermostOAuthenticator(app)
			t, err := oauth.GetToken(cc.ActingUserID)
			if err != nil {
				return nil, "", err
			}
			if t != nil {
				cc.ActingUserIsConnected = true
			}

			if expand.ActingUserAccessToken.Any() {
				if t != nil {
					call.Context.ActingUserAccessToken = t.AccessToken
				} else if expand.ActingUserAccessToken.IsRequired() {
					// If the Call requires OAuth token and we don't have one,
					// start OAuth2.
					connectURL, err := p.startMattermostOAuthConnect(oauth, cc.ActingUserID, app, call)
					return nil, connectURL, err
				}
			}
		}

		if expand.ActingUser != "" && e.ActingUser == nil {
			actingUser, err := p.mm.User.Get(e.ActingUserID)
			if err != nil {
				return nil, "", errors.Wrapf(err, "failed to expand acting user %s", e.ActingUserID)
			}
			e.ActingUser = actingUser
		}
	}

	if expand.Channel != "" && e.ChannelID != "" && e.Channel == nil {
		ch, err := p.mm.Channel.Get(e.ChannelID)
		if err != nil {
			return nil, "", errors.Wrapf(err, "failed to expand channel %s", e.ChannelID)
		}
		e.Channel = ch
	}

	// TODO expand Mentioned

	if expand.Post != "" && e.PostID != "" && e.Post == nil {
		post, err := p.mm.Post.GetPost(e.PostID)
		if err != nil {
			return nil, "", errors.Wrapf(err, "failed to expand post %s", e.PostID)
		}
		e.Post = post
	}

	if expand.RootPost != "" && e.RootPostID != "" && e.RootPost == nil {
		post, err := p.mm.Post.GetPost(e.RootPostID)
		if err != nil {
			return nil, "", errors.Wrapf(err, "failed to expand root post %s", e.RootPostID)
		}
		e.RootPost = post
	}

	if expand.Team != "" && e.TeamID != "" && e.Team == nil {
		team, err := p.mm.Team.Get(e.TeamID)
		if err != nil {
			return nil, "", errors.Wrapf(err, "failed to expand team %s", e.TeamID)
		}
		e.Team = team
	}

	if expand.User != "" && e.UserID != "" && e.User == nil {
		user, err := p.mm.User.Get(e.UserID)
		if err != nil {
			return nil, "", errors.Wrapf(err, "failed to expand user %s", e.UserID)
		}
		e.User = user
	}

	cc.ExpandedContext = api.ExpandedContext{
		BotAccessToken: app.BotAccessToken,

		ActingUser: stripUser(e.ActingUser, expand.ActingUser),
		App:        stripApp(app, expand.App),
		Channel:    stripChannel(e.Channel, expand.Channel),
		Post:       stripPost(e.Post, expand.Post),
		RootPost:   stripPost(e.RootPost, expand.RootPost),
		Team:       stripTeam(e.Team, expand.Team),
		User:       stripUser(e.User, expand.User),
		// TODO Mentioned
	}

	return &cc, "", nil
}

func stripUser(user *model.User, level api.ExpandLevel) *model.User {
	if user == nil || level == api.ExpandAll {
		return user
	}
	if level != api.ExpandSummary {
		return nil
	}
	return &model.User{
		BotDescription: user.BotDescription,
		DeleteAt:       user.DeleteAt,
		Email:          user.Email,
		FirstName:      user.FirstName,
		Id:             user.Id,
		IsBot:          user.IsBot,
		LastName:       user.LastName,
		Locale:         user.Locale,
		Nickname:       user.Nickname,
		Roles:          user.Roles,
		Timezone:       user.Timezone,
		Username:       user.Username,
	}
}

func stripChannel(channel *model.Channel, level api.ExpandLevel) *model.Channel {
	if channel == nil || level == api.ExpandAll {
		return channel
	}
	if level != api.ExpandSummary {
		return nil
	}
	return &model.Channel{
		Id:          channel.Id,
		DeleteAt:    channel.DeleteAt,
		TeamId:      channel.TeamId,
		Type:        channel.Type,
		DisplayName: channel.DisplayName,
		Name:        channel.Name,
	}
}

func stripTeam(team *model.Team, level api.ExpandLevel) *model.Team {
	if team == nil || level == api.ExpandAll {
		return team
	}
	if level != api.ExpandSummary {
		return nil
	}
	return &model.Team{
		Id:          team.Id,
		DisplayName: team.DisplayName,
		Name:        team.Name,
		Description: team.Description,
		Email:       team.Email,
		Type:        team.Type,
	}
}

func stripPost(post *model.Post, level api.ExpandLevel) *model.Post {
	if post == nil || level == api.ExpandAll {
		return post
	}
	if level != api.ExpandSummary {
		return nil
	}
	return &model.Post{
		Id:        post.Id,
		Type:      post.Type,
		UserId:    post.UserId,
		ChannelId: post.ChannelId,
		RootId:    post.RootId,
		Message:   post.Message,
	}
}

func stripApp(app *api.App, level api.ExpandLevel) *api.App {
	if app == nil {
		return nil
	}

	clone := *app
	clone.Secret = ""
	clone.OAuth2ClientSecret = ""

	switch level {
	case api.ExpandAll, api.ExpandSummary:
		return &clone
	}
	return nil
}
