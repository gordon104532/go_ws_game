package goWsGame

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/go-redis/redis/v8"
	"github.com/gordon104532/go_ws_game/controller"
	log "github.com/sirupsen/logrus"
)

type App struct {
	httpPort string
	cache    *redis.Client
}

func NewApp(
	httpPort string,
	cache *redis.Client,
) *App {
	return &App{
		httpPort: httpPort,
		cache:    cache,
	}
}

func (a *App) Run() {

	// twitchBotSrv := service.NewTwitchBotService()

	httpServer := controller.NewHttpServer(a.httpPort, a.cache)
	go httpServer.Run()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sc)

	<-sc
	log.Info("Service closing...")

	// err := twitchBotSrv.Stop()
	// if err != nil {
	// 	log.Errorf("twitchBot stop err: %s", err.Error())
	// }

	log.Info("Service closed")
}
