package otpService

// OtpConsumer consumer interface
type OtpConsumer interface {
	StartConsumer(workerPoolSize int, exchange, queueName, bindingKey, consumerTag string) error
}
