package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/apache/incubator-answer/plugin"
	"github.com/google/go-github/v50/github"
	"github.com/segmentfault/pacman/log"
	"golang.org/x/oauth2"
	oauth2Wechat "github.com/sanxia/goauth"
)

type Connector struct {
	Config *ConnectorConfig
}

type ConnectorConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func init() {
	plugin.Register(&Connector{
		Config: &ConnectorConfig{},
	})
}

func (g *Connector) Info() plugin.Info {
	return plugin.Info{
		Name:        "Wechat Oauth2",
		SlugName:    "wechat_connector",
		Description: "wechat oauth2 for answer connector plugin",
		Author:      "HassBox",
		Version:     "0.0.1",
		Link:        "https://github.com/apache/incubator-answer-plugins/tree/main/connector-github",
	}
}

func (g *Connector) ConnectorLogoSVG() string {
	return `PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIyNCIgaGVpZ2h0PSIyNCIgdmlld0JveD0iMCAwIDI0IDI0Ij48cGF0aCBkPSJNMTIgMGMtNi42MjYgMC0xMiA1LjM3My0xMiAxMiAwIDUuMzAyIDMuNDM4IDkuOCA4LjIwNyAxMS4zODcuNTk5LjExMS43OTMtLjI2MS43OTMtLjU3N3YtMi4yMzRjLTMuMzM4LjcyNi00LjAzMy0xLjQxNi00LjAzMy0xLjQxNi0uNTQ2LTEuMzg3LTEuMzMzLTEuNzU2LTEuMzMzLTEuNzU2LTEuMDg5LS43NDUuMDgzLS43MjkuMDgzLS43MjkgMS4yMDUuMDg0IDEuODM5IDEuMjM3IDEuODM5IDEuMjM3IDEuMDcgMS44MzQgMi44MDcgMS4zMDQgMy40OTIuOTk3LjEwNy0uNzc1LjQxOC0xLjMwNS43NjItMS42MDQtMi42NjUtLjMwNS01LjQ2Ny0xLjMzNC01LjQ2Ny01LjkzMSAwLTEuMzExLjQ2OS0yLjM4MSAxLjIzNi0zLjIyMS0uMTI0LS4zMDMtLjUzNS0xLjUyNC4xMTctMy4xNzYgMCAwIDEuMDA4LS4zMjIgMy4zMDEgMS4yMy45NTctLjI2NiAxLjk4My0uMzk5IDMuMDAzLS40MDQgMS4wMi4wMDUgMi4wNDcuMTM4IDMuMDA2LjQwNCAyLjI5MS0xLjU1MiAzLjI5Ny0xLjIzIDMuMjk3LTEuMjMuNjUzIDEuNjUzLjI0MiAyLjg3NC4xMTggMy4xNzYuNzcuODQgMS4yMzUgMS45MTEgMS4yMzUgMy4yMjEgMCA0LjYwOS0yLjgwNyA1LjYyNC01LjQ3OSA1LjkyMS40My4zNzIuODIzIDEuMTAyLjgyMyAyLjIyMnYzLjI5M2MwIC4zMTkuMTkyLjY5NC44MDEuNTc2IDQuNzY1LTEuNTg5IDguMTk5LTYuMDg2IDguMTk5LTExLjM4NiAwLTYuNjI3LTUuMzczLTEyLTEyLTEyeiIvPjwvc3ZnPg==`
}

func (g *Connector) ConnectorName() plugin.Translator {
	return "Wechat connector"
}

func (g *Connector) ConnectorSlugName() string {
	return "wechat"
}

func (g *Connector) ConnectorSender(ctx *plugin.GinContext, receiverURL string) (redirectURL string) {
	weChatOauth := oauth2Wechat.NewWeChat(g.Config.ClientID, g.Config.ClientSecret, receiverURL)
	return weChatOauth.GetAuthorizeUrl()
}

func (g *Connector) ConnectorReceiver(ctx *plugin.GinContext, receiverURL string) (userInfo plugin.ExternalLoginUserInfo, err error) {
	code := ctx.Query("code")
	// Exchange code for token
	weChatOauth := oauth2Wechat.NewWeChat(g.Config.ClientID, g.Config.ClientSecret, receiverURL)
	weChatToken, err := weChatOauth.GetAccessToken(code)
	if err != nil {
		return userInfo, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	weChatUserInfo, err := weChatOauth.GetUserInfo(weChatToken.AccessToken, weChatToken.OpenId)
	userInfo = plugin.ExternalLoginUserInfo{
		ExternalID:  "1",
		DisplayName: weChatUserInfo.Nickname,
		Username:    weChatUserInfo.Nickname,
		Avatar:      weChatUserInfo.Avatar,
		Email:       "",
		MetaInfo:    "",
	}
	return userInfo, nil
}

func (g *Connector) ConfigFields() []plugin.ConfigField {
	return []plugin.ConfigField{
		{
			Name:        "client_id",
			Type:        plugin.ConfigTypeInput,
			Title:       "client_id_1",
			Description: "client_id_2",
			Required:    true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			Value: g.Config.ClientID,
		},
		{
			Name:        "client_secret",
			Type:        plugin.ConfigTypeInput,
			Title:       "client_secret_1",
			Description: "client_secret_2",
			Required:    true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			Value: g.Config.ClientSecret,
		},
	}
}

func (g *Connector) ConfigReceiver(config []byte) error {
	c := &ConnectorConfig{}
	_ = json.Unmarshal(config, c)
	g.Config = c
	return nil
}