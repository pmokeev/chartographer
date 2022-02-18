package main

import (
	"context"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	server "pmokeev/chartographer/internal"
	"pmokeev/chartographer/internal/routers"
	"pmokeev/chartographer/internal/services"
	"syscall"
)

func initConfigFile() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func main() {
	if err := initConfigFile(); err != nil {
		log.Fatalf("Error while init config %s", err.Error())
	}

	service := services.NewService()
	chartRouter := routers.NewChartRouter(service)
	chartServer := server.NewServer()

	go func() {
		if err := chartServer.Run(viper.GetString("port"), chartRouter.InitChartRouter()); err != nil {
			log.Fatalf("Error while running server %s", err.Error())
		}
	}()

	log.Print("API started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	log.Print("API shutdowned")

	if err := chartServer.Shutdown(context.Background()); err != nil {
		log.Fatalf("Error while shutdowning server %s", err.Error())
	}
}
