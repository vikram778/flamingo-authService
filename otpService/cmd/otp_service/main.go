package main

import (
	"flamingo-authService/config"
	"flamingo-authService/otpService/server"
	"flamingo-authService/pkg/log"
	"flamingo-authService/pkg/rabbitmq"
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

	s := server.NewOtpServer(amqpConn, cfg)

	log.Info("Starting OTP Service")

	s.Run()
}
