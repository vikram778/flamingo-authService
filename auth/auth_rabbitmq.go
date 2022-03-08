package auth

type Publisher interface {
	Publish(body []byte) error
}
