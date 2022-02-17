package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
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

	fmt.Println("Goodbye, cruel world!")
}
