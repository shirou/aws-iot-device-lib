package mqttutils

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Unsubscribe is a utility function about unsubscribing topics
func Unsubscribe(cli mqtt.Client, topics []string) error {
	token := cli.Unsubscribe(topics...)
	token.Wait()
	return token.Error()
}

// Subscribe is a utility function about subscribing topics
func Subscribe(cli mqtt.Client, topics []string, qos int, callback mqtt.MessageHandler) error {
	filter := make(map[string]byte)
	for _, t := range topics {
		filter[t] = byte(qos)
	}

	token := cli.SubscribeMultiple(filter, callback)
	token.Wait()
	return token.Error()
}

// Publish is a utility function about subscribing topics
// Note: retain always false
func Publish(cli mqtt.Client, topic string, qos int, payload []byte) error {
	token := cli.Publish(topic, byte(qos), false, payload)
	token.Wait()
	return token.Error()
}
