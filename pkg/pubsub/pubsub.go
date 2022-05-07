package pubsub

type PubSub interface {
	Publish(topic, message string) error
	Subscribe(topic string, callback func(string) error) error
}
