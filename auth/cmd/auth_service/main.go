package main

import (
	"flamingo-authService/auth/server"
	"flamingo-authService/config"
	"flamingo-authService/pkg/log"
	"flamingo-authService/pkg/postgres"
	"flamingo-authService/pkg/rabbitmq"
	"fmt"
	"go.uber.org/zap"

	"os"
)

func main() {

	log.SetLogLevel()

	log.Info("Starting server")

	configPath := config.GetConfigPath(os.Getenv("config"))
	cfg, err := config.GetConfig(configPath)
	if err != nil {
		log.Fatal("Loading config:", zap.Error(err))
	}

	amqpConn, err := rabbitmq.NewRabbitMQConn(cfg)
	if err != nil {
		log.Fatal("RabbitMQ Error", zap.Error(err))
	}
	defer amqpConn.Close()

	psqlDB, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		log.Fatal(fmt.Sprintf("Postgresql init: %s", err))
	}
	defer psqlDB.Close()

	log.Info(fmt.Sprintf("PostgreSQL connected: %#v", psqlDB.Stats()))

	s := server.NewAuthServer(amqpConn, cfg, psqlDB)

	log.Fatal("Starting Service Error", zap.Error(s.Run()))
}
