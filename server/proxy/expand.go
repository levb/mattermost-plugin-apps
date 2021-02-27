package proxy

import (
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/mattermost/mattermost-plugin-apps/apps"
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

func (p *proxy) expandCall(inCall *apps.Call, app *apps.App, adminAccessToken string, oauth oauther.OAuther, cache *expandCache) (*apps.Call, error) {
	call := *inCall
	if cache == nil {
		cache = &expandCache{}
	}
	cc := apps.Context{}
	if call.Context != nil {
		cc = *call.Context
	}
	expand := apps.Expand{}
	if call.Expand != nil {
		expand = *call.Expand
	}

	conf := p.conf.Get()
	cc.MattermostSiteURL = conf.MattermostSiteURL

	cc.BotUserID = app.BotUserID
	if app.GrantedPermissions.Contains(apps.PermissionActAsBot) {
		cc.ExpandedContext.BotAccessToken = app.BotAccessToken
	}

	if cc.ActingUserID != "" {
		if app.GrantedPermissions.Contains(apps.PermissionActAsUser) {
			t, err := oauth.GetToken(cc.ActingUserID)
			if err != nil {
				return nil, err
			}
			if t != nil {
				cc.ActingUserIsConnected = true
			}

			if expand.ActingUserAccessToken.Any() {
				if t != nil {
					cc.ExpandedContext.ActingUserAccessToken = t.AccessToken
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
		cc.ExpandedContext.AdminAccessToken = adminAccessToken
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

	cc.ExpandedContext.BotAccessToken = app.BotAccessToken
	cc.ExpandedContext.ActingUser = stripUser(cache.actingUser, expand.ActingUser)
	cc.ExpandedContext.App = stripApp(app, expand.App)
	cc.ExpandedContext.Channel = stripChannel(cache.channel, expand.Channel)
	cc.ExpandedContext.Post = stripPost(cache.post, expand.Post)
	cc.ExpandedContext.RootPost = stripPost(cache.rootPost, expand.RootPost)
	cc.ExpandedContext.Team = stripTeam(cache.team, expand.Team)
	cc.ExpandedContext.User = stripUser(cache.user, expand.User)
	// TODO Mentioned

	call.Context = &cc
	return &call, nil
}

func stripUser(user *model.User, level apps.ExpandLevel) *model.User {
	if user == nil || level == apps.ExpandAll {
		return user
	}
	if level != apps.ExpandSummary {
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

func stripChannel(channel *model.Channel, level apps.ExpandLevel) *model.Channel {
	if channel == nil || level == apps.ExpandAll {
		return channel
	}
	if level != apps.ExpandSummary {
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

func stripTeam(team *model.Team, level apps.ExpandLevel) *model.Team {
	if team == nil || level == apps.ExpandAll {
		return team
	}
	if level != apps.ExpandSummary {
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

func stripPost(post *model.Post, level apps.ExpandLevel) *model.Post {
	if post == nil || level == apps.ExpandAll {
		return post
	}
	if level != apps.ExpandSummary {
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

func stripApp(app *apps.App, level apps.ExpandLevel) *apps.App {
	if app == nil {
		return nil
	}

	clone := *app
	clone.Secret = ""
	clone.OAuth2ClientSecret = ""

	switch level {
	case apps.ExpandAll, apps.ExpandSummary:
		return &clone
	}
	return nil
}
