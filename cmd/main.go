package main

import (
	"context"
	"errors"
	server "github.com/pmokeev/chartographer/internal"
	"github.com/pmokeev/chartographer/internal/routers"
	"github.com/pmokeev/chartographer/internal/services"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
		if err := chartServer.Run(viper.GetString("port"), chartRouter.InitChartRouter()); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("Listen: %s\n", err)
		}
	}()

	log.Println("API started")

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down API...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := chartServer.Shutdown(ctx); err != nil {
		log.Fatal("API forced to shutdown:", err)
	}

	log.Println("API exiting")
}
