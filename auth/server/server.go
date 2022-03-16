package server

import (
	authGrpc "flamingo-authService/auth/grpc"
	authService "flamingo-authService/auth/proto"
	"flamingo-authService/auth/rabbitmq"
	"flamingo-authService/auth/repository"
	"flamingo-authService/config"
	"flamingo-authService/pkg/log"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"os/signal"
	"syscall"
)

// Server struct
type Server struct {
	db       *sqlx.DB
	amqpConn *amqp.Connection
	cfg      *config.Config
}

// NewAuthServer constructor
func NewAuthServer(amqpConn *amqp.Connection, cfg *config.Config, db *sqlx.DB) *Server {
	return &Server{amqpConn: amqpConn, cfg: cfg, db: db}
}

// Run server
func (s *Server) Run() error {

	authPublisher, err := rabbitmq.NewPublisher(s.cfg, s.amqpConn)
	if err != nil {
		return err
	}
	defer authPublisher.CloseChan()
	log.Info("Auth Publisher initialized")

	authRepository := repository.NewDBOpsRepository(s.db)

	l, err := net.Listen("tcp", s.cfg.Server.Port)
	if err != nil {
		return err
	}
	defer func(l net.Listener) {
		_ = l.Close()
	}(l)

	err = authPublisher.SetupExchangeAndQueue(s.cfg.RabbitMQ.Exchange, s.cfg.RabbitMQ.Queue, s.cfg.RabbitMQ.RoutingKey)
	if err != nil {
		return err
	}
	server := grpc.NewServer()

	authGrpcMicroservice := authGrpc.NewAuthMicroservice(*authRepository, s.cfg, *authPublisher)
	authService.RegisterAuthServiceServer(server, authGrpcMicroservice)

	log.Info("Auth Service initialized")

	if s.cfg.Server.Mode != "Production" {
		reflection.Register(server)
	}

	go func() {
		log.Info(fmt.Sprintf("Server is listening on port: %v", s.cfg.Server.Port))
		log.Fatal("Error starting Service", zap.Error(server.Serve(l)))
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		log.Error(fmt.Sprintf("signal.Notify: %v", v))
	}

	server.GracefulStop()
	log.Info("Server Exited Properly")

	return nil
}
