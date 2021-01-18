package proxy

import (
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/mattermost/mattermost-plugin-apps/server/api"
	"github.com/mattermost/mattermost-plugin-apps/server/utils/oauther"
)

type expandCache struct {
	actingUser *model.User
	channel    *model.Channel
	post       *model.Post
	rootPost   *model.Post
	team       *model.Team
	user       *model.User
}

var errOAuthRequired = errors.New("oauth2 (re-)authorization with Mattermost required")

func (p *Proxy) expandCall(inCall *api.Call, app *api.App, sessionToken api.SessionToken, oauth oauther.OAuther, cache *expandCache) (*api.Call, error) {
	call := *inCall
	if cache == nil {
		cache = &expandCache{}
	}
	cc := api.Context{}
	if call.Context != nil {
		cc = *call.Context
	}
	expand := api.Expand{}
	if call.Expand != nil {
		expand = *call.Expand
	}

	conf := p.conf.GetConfig()
	cc.MattermostSiteURL = conf.MattermostSiteURL

	cc.BotUserID = app.BotUserID
	if app.GrantedPermissions.Contains(api.PermissionActAsBot) {
		cc.BotAccessToken = app.BotAccessToken
	}

	if cc.ActingUserID != "" {
		if app.GrantedPermissions.Contains(api.PermissionActAsUser) {
			t, err := oauth.GetToken(cc.ActingUserID)
			if err != nil {
				return nil, err
			}
			if t != nil {
				cc.ActingUserIsConnected = true
			}

			if expand.ActingUserAccessToken.Any() {
				if t != nil {
					cc.ActingUserAccessToken = t.AccessToken
				} else if expand.ActingUserAccessToken.IsRequired() {
					return nil, errOAuthRequired
				}
			}
		}

		if expand.ActingUser != "" && cache.actingUser == nil {
			actingUser, err := p.mm.User.Get(cc.ActingUserID)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to expand acting user %s", cc.ActingUserID)
			}
			cache.actingUser = actingUser
		}
	}

	// TODO: do we need another way of obtaining an admin token? Should it be
	// out of expand?
	// TODO: Implement collecting user consent for the admin token, in-line?
	if expand.AdminAccessToken.Any() {
		cc.AdminAccessToken = string(sessionToken)
	}

	if expand.Channel != "" && cc.ChannelID != "" && cache.channel == nil {
		ch, err := p.mm.Channel.Get(cc.ChannelID)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to expand channel %s", cc.ChannelID)
		}
		cache.channel = ch
	}

	if expand.Post != "" && cc.PostID != "" && cache.post == nil {
		post, err := p.mm.Post.GetPost(cc.PostID)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to expand post %s", cc.PostID)
		}
		cache.post = post
	}

	if expand.RootPost != "" && cc.RootPostID != "" && cache.rootPost == nil {
		post, err := p.mm.Post.GetPost(cc.RootPostID)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to expand root post %s", cc.RootPostID)
		}
		cache.rootPost = post
	}

	if expand.Team != "" && cc.TeamID != "" && cache.team == nil {
		team, err := p.mm.Team.Get(cc.TeamID)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to expand team %s", cc.TeamID)
		}
		cache.team = team
	}

	// TODO expand Mentioned
	if expand.User != "" && cc.UserID != "" && cache.user == nil {
		user, err := p.mm.User.Get(cc.UserID)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to expand user %s", cc.UserID)
		}
		cache.user = user
	}

	cc.ExpandedContext = api.ExpandedContext{
		BotAccessToken: app.BotAccessToken,

		ActingUser: stripUser(cache.actingUser, expand.ActingUser),
		App:        stripApp(app, expand.App),
		Channel:    stripChannel(cache.channel, expand.Channel),
		Post:       stripPost(cache.post, expand.Post),
		RootPost:   stripPost(cache.rootPost, expand.RootPost),
		Team:       stripTeam(cache.team, expand.Team),
		User:       stripUser(cache.user, expand.User),
		// TODO Mentioned
	}

	call.Context = &cc
	return &call, nil
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
