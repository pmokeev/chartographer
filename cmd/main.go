package main

import (
	server "github.com/pmokeev/chartographer/internal"
	"github.com/pmokeev/chartographer/internal/routers"
	"github.com/pmokeev/chartographer/internal/services"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
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

	service := services.NewService(os.Args[1])
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
}
