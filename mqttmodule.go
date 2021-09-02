package mqttmodule

import (
	"errors"
)

var subscribedChannels map[string]chan MqttMessage = make(map[string]chan MqttMessage)

func Subscribe(topics []string, channels []chan MqttMessage) error {

	if len(topics) != len(channels) {
		return errors.New("length of channels must be equal to the length of topics")
	}
	conn := connectMqtt()

	for i := range topics {
		err := subscribe(conn, topics[i], channels[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func Publish(pm chan MqttMessage) {
	conn := connectMqtt()

	for {
		msg := <-pm
		conn.Publish(msg.Topic, 0, false, []byte(msg.Data))
	}
}
