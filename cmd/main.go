package main

import (
	"github.com/pmokeev/chartographer/internal/routers"
	"github.com/pmokeev/chartographer/internal/services"
	"github.com/spf13/viper"
	"log"
	"os"
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

	if err := chartRouter.InitChartRouter().Run(":8000"); err != nil {
		log.Fatalf("Error while running server %s", err.Error())
	}
}
