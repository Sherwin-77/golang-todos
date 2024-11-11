package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4/middleware"
	"github.com/sherwin-77/go-echo-template/configs"
	"github.com/sherwin-77/go-echo-template/internal/builder"
	"github.com/sherwin-77/go-echo-template/internal/http/handler"
	"github.com/sherwin-77/go-echo-template/pkg/caches"
	"github.com/sherwin-77/go-echo-template/pkg/database"
	"github.com/sherwin-77/go-echo-template/pkg/server"
)

func main() {
	config := configs.LoadConfig()

	db, err := database.InitDB(config.Postgres)
	if err != nil {
		panic(err)
	}

	cache := caches.NewCache(caches.InitRedis(config.Redis))

	echoServer := server.NewServer()
	echoServer.Use(middleware.LoggerWithConfig(configs.GetEchoLoggerConfig()))
	echoServer.Use(middleware.RecoverWithConfig(configs.GetEchoRecoverConfig()))
	echoServer.Validator = configs.NewAppValidator()
	echoServer.HTTPErrorHandler = handler.HTTPErrorHandler

	group := echoServer.Group("/api")
	builder.BuildV1Routes(config, db, cache, group)

	runServer(echoServer, config)
	waitForShutdown(echoServer)
}

func runServer(s *server.Server, config *configs.Config) {
	go func() {
		if err := s.Start(fmt.Sprintf("localhost:%s", config.Port)); err != nil {
			s.Logger.Fatal(err)
		}
	}()
}

func waitForShutdown(s *server.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go func() {
		if err := s.Shutdown(ctx); err != nil {
			s.Logger.Fatal(err)
		}
	}()
}
