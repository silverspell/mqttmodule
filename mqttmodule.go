package mqttmodule

import (
	"errors"
)

var subscribedChannels map[string]chan []byte = make(map[string]chan []byte)

func Subscribe(topics []string, channels []chan []byte) error {

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

func Publish(pm chan PublishMessage) {
	conn := connectMqtt()

	for {
		msg := <-pm
		conn.Publish(msg.Topic, 0, false, []byte(msg.Data))
	}
}
