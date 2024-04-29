package service

import (
	"context"
	"os"

	twitch "github.com/gempir/go-twitch-irc/v3"
	"github.com/gordon104532/go_ws_game/pkg"
	log "github.com/sirupsen/logrus"
)

type TwitchBotService struct {
	client           *twitch.Client
	retryCount       int64
	subscribeChannel []string
}

func NewTwitchBotService() *TwitchBotService {
	log.Info("TwitchBotService init...")

	service := &TwitchBotService{
		subscribeChannel: []string{"gordon1045321"},
	}

	// 聊天室用oauthi
	var twitchToken string
	var err error
	if os.Getenv("ENV") == "docker" {
		twitchToken = os.Getenv("TWITCH_TOKEN")
		log.Info("TWITCH_TOKEN: ", twitchToken)
	} else {
		twitchToken, err = pkg.GetSecret(context.Background(), "projects/964403411847/secrets/twitch_oauth_token/versions/latest")
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	service.client = twitch.NewClient("gordon1045321", twitchToken)

	service.client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		service.msgHandle(message)
	})

	// 成功連線
	service.client.OnConnect(func() {
		log.Info("TwitchBot OnConnect")
		service.retryCount = 0
	})

	// 重連事件
	service.client.OnReconnectMessage(func(message twitch.ReconnectMessage) {
		log.Info("[OnReconnectMessage]", message)
	})

	// 加入頻道
	service.client.Join(service.subscribeChannel...)

	go service.Run()

	return service
}

func (t *TwitchBotService) Run() error {
	log.Info("TwitchBot Start")
	err := t.client.Connect()
	if err != nil {
		log.Fatal("TwitchBot Connect error: ", err)
		return err
	}

	return nil
}

func (t *TwitchBotService) Stop() error {
	err := t.client.Disconnect()
	if err != nil {
		return err
	}

	return nil
}

func (t *TwitchBotService) msgHandle(message twitch.PrivateMessage) {
	log.Info("[" + message.Channel + "]" + message.User.DisplayName + ": " + message.Message)
}
