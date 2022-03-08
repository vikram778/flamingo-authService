package rabbitmq

import (
	"flamingo-authService/config"
	"flamingo-authService/pkg/rabbitmq"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"time"
)

const (
	exchangeKind       = "direct"
	exchangeDurable    = true
	exchangeAutoDelete = false
	exchangeInternal   = false
	exchangeNoWait     = false

	queueDurable    = true
	queueAutoDelete = false
	queueExclusive  = false
	queueNoWait     = false

	publishMandatory = false
	publishImmediate = false
)

type Publisher struct {
	amqpChan *amqp.Channel
	cfg      *config.Config
}

// NewPublisher Emails rabbitmq publisher constructor
func NewPublisher(cfg *config.Config) (*Publisher, error) {
	mqConn, err := rabbitmq.NewRabbitMQConn(cfg)
	if err != nil {
		return nil, err
	}
	amqpChan, err := mqConn.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "p.amqpConn.Channel")
	}

	return &Publisher{cfg: cfg, amqpChan: amqpChan}, nil
}

// SetupExchangeAndQueue create exchange and queue
func (p *Publisher) SetupExchangeAndQueue(exchange, queueName, bindingKey, consumerTag string) error {

	err := p.amqpChan.ExchangeDeclare(
		exchange,
		exchangeKind,
		exchangeDurable,
		exchangeAutoDelete,
		exchangeInternal,
		exchangeNoWait,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "Error ch.ExchangeDeclare")
	}

	queue, err := p.amqpChan.QueueDeclare(
		queueName,
		queueDurable,
		queueAutoDelete,
		queueExclusive,
		queueNoWait,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "Error ch.QueueDeclare")
	}

	err = p.amqpChan.QueueBind(
		queue.Name,
		bindingKey,
		exchange,
		queueNoWait,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "Error ch.QueueBind")
	}

	return nil
}

// CloseChan Close messages chan
func (p *Publisher) CloseChan() {
	if err := p.amqpChan.Close(); err != nil {
	}
}

// Publish message
func (p *Publisher) Publish(body []byte) error {

	if err := p.amqpChan.Publish(
		p.cfg.RabbitMQ.Exchange,
		p.cfg.RabbitMQ.RoutingKey,
		publishMandatory,
		publishImmediate,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Body:         body,
		},
	); err != nil {
		return errors.Wrap(err, "ch.Publish")
	}

	return nil
}
