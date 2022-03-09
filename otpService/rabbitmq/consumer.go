package rabbitmq

import (
	"context"
	"flamingo-authService/otpService/otpOps"
	"flamingo-authService/pkg/log"
	"fmt"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
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

	prefetchCount  = 1
	prefetchSize   = 0
	prefetchGlobal = false

	consumeAutoAck   = false
	consumeExclusive = false
	consumeNoLocal   = false
	consumeNoWait    = false
)

// OtpConsumer Images Rabbitmq consumer
type OtpConsumer struct {
	amqpConn *amqp.Connection
	otpOps   otpOps.OtpOps
}

// NewOtpConsumer Consumer constructor
func NewOtpConsumer(amqpConn *amqp.Connection, otpOps otpOps.OtpOps) *OtpConsumer {
	return &OtpConsumer{amqpConn: amqpConn, otpOps: otpOps}
}

// CreateChannel Consume messages
func (c *OtpConsumer) CreateChannel(exchangeName, queueName, bindingKey, consumerTag string) (*amqp.Channel, error) {
	ch, err := c.amqpConn.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "Error amqpConn.Channel")
	}

	log.Info(fmt.Sprintf("Declaring exchange: %s", exchangeName))
	err = ch.ExchangeDeclare(
		exchangeName,
		exchangeKind,
		exchangeDurable,
		exchangeAutoDelete,
		exchangeInternal,
		exchangeNoWait,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error ch.ExchangeDeclare")
	}

	queue, err := ch.QueueDeclare(
		queueName,
		queueDurable,
		queueAutoDelete,
		queueExclusive,
		queueNoWait,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error ch.QueueDeclare")
	}

	log.Info(fmt.Sprintf("Declared queue, binding it to exchange: Queue: %v, messagesCount: %v, "+
		"consumerCount: %v, exchange: %v, bindingKey: %v",
		queue.Name,
		queue.Messages,
		queue.Consumers,
		exchangeName,
		bindingKey,
	))

	err = ch.QueueBind(
		queue.Name,
		bindingKey,
		exchangeName,
		queueNoWait,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error ch.QueueBind")
	}

	log.Info(fmt.Sprintf("Queue bound to exchange, starting to consume from queue, consumerTag: %v", consumerTag))

	err = ch.Qos(
		prefetchCount,  // prefetch count
		prefetchSize,   // prefetch size
		prefetchGlobal, // global
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error  ch.Qos")
	}

	return ch, nil
}

func (c *OtpConsumer) worker(ctx context.Context, messages <-chan amqp.Delivery) {

	for delivery := range messages {
		log.Info(fmt.Sprintf("processDeliveries deliveryTag %v", delivery.DeliveryTag))

		err := c.otpOps.SendOtp(ctx, delivery.Body)
		if err != nil {
			if err := delivery.Reject(false); err != nil {
				log.Error("Err delivery.Reject: %v", zap.Error(err))
			}
			log.Error("Failed to process delivery: %v", zap.Error(err))
		} else {
			err = delivery.Ack(false)
			if err != nil {
				log.Error("Failed to acknowledge delivery: %v", zap.Error(err))
			}
		}
	}

	log.Info("Deliveries channel closed")
}

// StartConsumer Start new rabbitmq consumer
func (c *OtpConsumer) StartConsumer(workerPoolSize int, exchange, queueName, bindingKey, consumerTag string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, err := c.CreateChannel(exchange, queueName, bindingKey, consumerTag)
	if err != nil {
		return errors.Wrap(err, "CreateChannel")
	}
	defer ch.Close()

	deliveries, err := ch.Consume(
		queueName,
		consumerTag,
		consumeAutoAck,
		consumeExclusive,
		consumeNoLocal,
		consumeNoWait,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "Consume")
	}

	for i := 0; i < workerPoolSize; i++ {
		go c.worker(ctx, deliveries)
	}

	chanErr := <-ch.NotifyClose(make(chan *amqp.Error))
	log.Error("ch.NotifyClose:", zap.Error(chanErr))
	return chanErr
}
