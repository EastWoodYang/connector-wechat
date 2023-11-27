package wechat

import (
	"encoding/json"
	"fmt"

	"github.com/apache/incubator-answer-plugins/connector-github/i18n"
	"github.com/apache/incubator-answer-plugins/connector-github/i18n"
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
	return `<svg enable-background="new 0 0 2499.7 2024.2" viewBox="0 0 2499.7 2024.2" xmlns="http://www.w3.org/2000/svg"><path d="m2499.7 1313.8c0-347.3-335.9-630.6-749.5-630.6s-747.2 283.3-747.2 630.6 335.9 630.6 749.5 630.6c80 0 155.4-9.1 226.2-29.7 20.6-4.6 43.4-2.3 64 6.9l185.1 100.5c11.4 6.9 27.4-4.6 22.8-18.3l-36.6-148.5c-4.6-22.8 2.3-45.7 22.8-59.4 160.1-116.5 262.9-287.9 262.9-482.1zm-1016.8-82.2c-57.1 0-102.8-45.7-102.8-102.8s45.7-102.8 102.8-102.8 102.8 45.7 102.8 102.8-45.7 102.8-102.8 102.8zm505 0c-57.1 0-102.8-45.7-102.8-102.8s45.7-102.8 102.8-102.8 102.8 45.7 102.8 102.8-48 102.8-102.8 102.8z"/><path d="m941.4 1316.1c0-386.1 365.6-699.2 818-699.2h34.3c-73.2-349.6-443.3-616.9-888.9-616.9-500.4 0-904.8 333.6-904.8 747.2 0 228.5 125.7 434.1 322.2 571.2 13.7 9.1 18.3 25.1 16 41.1l-68.5 242.2c-4.6 16 13.7 32 29.7 22.8l267.3-159.9c18.3-11.4 38.8-13.7 59.4-6.9 89.1 22.8 182.8 36.6 278.8 36.6 20.6 0 41.1 0 64-2.3-18.4-57.1-27.5-116.5-27.5-175.9zm267.3-920.8c66.3 0 121.1 54.8 121.1 121.1s-54.8 121.1-121.1 121.1-121.1-54.8-121.1-121.1 54.9-121.1 121.1-121.1zm-607.8 242.2c-66.3 0-121.1-54.8-121.1-121.1s54.8-121.1 121.1-121.1 121.1 54.8 121.1 121.1-52.5 121.1-121.1 121.1z"/></svg>`
}

func (g *Connector) ConnectorName() plugin.Translator {
	return plugin.MakeTranslator(i18n.ConnectorName)
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