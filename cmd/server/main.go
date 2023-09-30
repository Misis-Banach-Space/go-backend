package main

import (
	"fmt"

	"github.com/yogenyslav/kokoc-hack/internal/api/router"
	"github.com/yogenyslav/kokoc-hack/internal/config"
	"github.com/yogenyslav/kokoc-hack/internal/database"
	"github.com/yogenyslav/kokoc-hack/internal/logging"
	"github.com/yogenyslav/kokoc-hack/internal/service"
)

func main() {
	if err := config.NewConfig(); err != nil {
		fmt.Printf("failed to init config: %+v", err)
		panic(err)
	}

	if err := logging.NewLogger(); err != nil {
		fmt.Printf("failed to init logger: %+v", err)
		panic(err)
	}

	pg, err := database.NewPostgres(10)
	if err != nil {
		logging.Log.Panicf("failed to create db instance: %+v", err)
	}
	defer pg.Close()
	logging.Log.Infof("initialised db instance")

	rabbitmq, err := service.NewRabbutMQ(pg.GetPool())
	if err != nil {
		logging.Log.Fatalf("failed to create rabbitmq instance: %+v", err)
	}
	defer rabbitmq.Close()

	r := router.NewRouter(pg, rabbitmq)
	logging.Log.Infof("starting server on port %s", config.Cfg.ServerPort)
	if err := r.Run(":" + config.Cfg.ServerPort); err != nil {
		logging.Log.Fatalf("failed to start server: %+v", err)
	}
}
