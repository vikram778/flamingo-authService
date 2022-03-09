package server

import (
	"context"
	"flamingo-authService/config"
	"flamingo-authService/otpService/otpOps"
	"flamingo-authService/otpService/rabbitmq"
	"flamingo-authService/pkg/log"
	"fmt"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Server struct
type Server struct {
	amqpConn *amqp.Connection
	cfg      *config.Config
}

// NewOtpServer constructor
func NewOtpServer(amqpConn *amqp.Connection, cfg *config.Config) *Server {
	return &Server{amqpConn: amqpConn, cfg: cfg}
}

func (s *Server) Run() {

	var wg sync.WaitGroup

	otpops := otpOps.NewOtpOps(s.cfg)
	otpConsumer := rabbitmq.NewOtpConsumer(s.amqpConn, *otpops)

	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)

	go func() {
		err := otpConsumer.StartConsumer(
			s.cfg.RabbitMQ.WorkerPoolSize,
			s.cfg.RabbitMQ.Exchange,
			s.cfg.RabbitMQ.Queue,
			s.cfg.RabbitMQ.RoutingKey,
			s.cfg.RabbitMQ.ConsumerTag,
		)
		if err != nil {
			log.Error("StartConsumer: %v", zap.Error(err))
			cancel()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		log.Error(fmt.Sprintf("signal.Notify: %v", v))
	case done := <-ctx.Done():
		log.Error(fmt.Sprintf("ctx.Done: %v", done))
	}

	wg.Wait()
}
