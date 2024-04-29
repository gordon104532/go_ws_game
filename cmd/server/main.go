package main

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	app "github.com/gordon104532/go_ws_game"
	"github.com/gordon104532/go_ws_game/cnf"
	"github.com/gordon104532/go_ws_game/pkg"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Service start...")

	cnf.InitConfig(cnf.Get())

	// cache
	cache, err := pkg.NewRedisConn(
		&redis.Options{
			Addr:     cnf.Get().Redis.Addr,
			Username: cnf.Get().Redis.Username,
			Password: cnf.Get().Redis.Password,
			DB:       cnf.Get().Redis.DB,
		},
	)
	if err != nil {
		log.Fatalf("Redis connection error: %s", err.Error())
	}
	app.NewApp(
		fmt.Sprintf("%d", cnf.Get().Server.Http.Port),
		cache,
	).Run()
}
