package wechat

import (
	"encoding/json"
	"fmt"

	"github.com/EastWoodYang/connector-wechat/i18n"
	"github.com/apache/incubator-answer/plugin"
	oauth2Wechat "github.com/EastWoodYang/goauth"
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
		Name:        plugin.MakeTranslator(i18n.InfoName),
		SlugName:    "wechat_connector",
		Description: plugin.MakeTranslator(i18n.InfoDescription),
		Author:      "HassBox",
		Version:     "0.0.1",
		Link:        "https://github.com/EastWoodYang/connector-wechat",
	}
}

func (g *Connector) ConnectorLogoSVG() string {
	return `PHN2ZyB3aWR0aD0iMjQiIGhlaWdodD0iMjQiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIgZmlsbC1ydWxlPSJldmVub2RkIiBjbGlwLXJ1bGU9ImV2ZW5vZGQiPjxwYXRoIGQ9Ik0xOSAyNGgtMTRjLTIuNzYxIDAtNS0yLjIzOS01LTV2LTE0YzAtMi43NjEgMi4yMzktNSA1LTVoMTRjMi43NjIgMCA1IDIuMjM5IDUgNXYxNGMwIDIuNzYxLTIuMjM4IDUtNSA1em0tLjY2NS02Ljk4NGMxLjAxNi0uNzM2IDEuNjY1LTEuODI1IDEuNjY1LTMuMDM1IDAtMi4yMTgtMi4xNTgtNC4wMTYtNC44MTktNC4wMTYtMi42NjIgMC00LjgxOSAxLjc5OC00LjgxOSA0LjAxNnMyLjE1NyA0LjAxNiA0LjgxOSA0LjAxNmMuNTUgMCAxLjA4MS0uMDc5IDEuNTczLS4yMmwuMTQxLS4wMjFjLjA5MyAwIC4xNzcuMDI4LjI1Ni4wNzRsMS4wNTUuNjA5LjA5My4wM2MuMDg5IDAgLjE2MS0uMDcyLjE2MS0uMTYxbC0uMDI2LS4xMTctLjIxNy0uODExLS4wMTctLjEwMmMwLS4xMDguMDUzLS4yMDMuMTM1LS4yNjJ6bS04LjU1Mi0xMS40ODVjLTMuMTk0IDAtNS43ODMgMi4xNTgtNS43ODMgNC44MiAwIDEuNDUyLjc3OSAyLjc1OSAxLjk5OCAzLjY0Mi4wOTguMDcuMTYyLjE4NS4xNjIuMzE0bC0uMDIuMTIzLS4yNjEuOTcyLS4wMzEuMTQxYzAgLjEwNy4wODYuMTkzLjE5My4xOTNsLjExMS0uMDM2IDEuMjY2LS43MzFjLjA5Ni0uMDU1LjE5Ni0uMDg5LjMwNy0uMDg5bC4xNy4wMjVjLjU5MS4xNyAxLjIyOC4yNjUgMS44ODguMjY1bC4zMTgtLjAwOGMtLjEyNi0uMzc2LS4xOTQtLjc3Mi0uMTk0LTEuMTgxIDAtMi40MjcgMi4zNjEtNC4zOTUgNS4yNzQtNC4zOTVsLjMxNC4wMDhjLS40MzYtMi4zMDItMi44MjctNC4wNjMtNS43MTItNC4wNjN6bTMuNzkxIDcuODA3Yy0uMzU1IDAtLjY0Mi0uMjg3LS42NDItLjY0MiAwLS4zNTUuMjg3LS42NDMuNjQyLS42NDMuMzU1IDAgLjY0My4yODguNjQzLjY0MyAwIC4zNTUtLjI4OC42NDItLjY0My42NDJ6bTMuMjEzIDBjLS4zNTUgMC0uNjQyLS4yODctLjY0Mi0uNjQyIDAtLjM1NS4yODctLjY0My42NDItLjY0My4zNTUgMCAuNjQzLjI4OC42NDMuNjQzIDAgLjM1NS0uMjg4LjY0Mi0uNjQzLjY0MnptLTguOTMyLTMuNzU5Yy0uNDI1IDAtLjc3MS0uMzQ1LS43NzEtLjc3MSAwLS40MjYuMzQ2LS43NzEuNzcxLS43NzEuNDI2IDAgLjc3Mi4zNDUuNzcyLjc3MSAwIC40MjYtLjM0Ni43NzEtLjc3Mi43NzF6bTMuODU2IDBjLS40MjYgMC0uNzcxLS4zNDUtLjc3MS0uNzcxIDAtLjQyNi4zNDUtLjc3MS43NzEtLjc3MS40MjYgMCAuNzcxLjM0NS43NzEuNzcxIDAgLjQyNi0uMzQ1Ljc3MS0uNzcxLjc3MXoiLz48L3N2Zz4=`
}

func (g *Connector) ConnectorName() plugin.Translator {
	return plugin.MakeTranslator(i18n.ConnectorName)
}

func (g *Connector) ConnectorSlugName() string {
	return "wechat_connector"
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
		ExternalID:  weChatToken.UnionId,
		DisplayName: weChatUserInfo.Nickname,
		Username:    weChatUserInfo.Nickname,
		Avatar:      weChatUserInfo.Avatar,
		Email:       "temp_" + weChatToken.UnionId + "@hassbox.cn",
		MetaInfo:    "",
	}
	return userInfo, nil
}

func (g *Connector) ConfigFields() []plugin.ConfigField {
	return []plugin.ConfigField{
		{
			Name:        "client_id",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigClientIDTitle),
			Description: plugin.MakeTranslator(i18n.ConfigClientIDDescription),
			Required:    true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			Value: g.Config.ClientID,
		},
		{
			Name:        "client_secret",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigClientSecretTitle),
			Description: plugin.MakeTranslator(i18n.ConfigClientSecretDescription),
			Required:    true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypePassword,
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