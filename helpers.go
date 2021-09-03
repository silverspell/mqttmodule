package mqttmodule

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func onConnectHandler(client mqtt.Client) {
	for k, v := range subscribedChannels {
		err := subscribe(client, k, v)
		failOnError(err, "Error on subscription")
	}
}

func onMessageHandler(client mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	if channel, ok := subscribedChannels[topic]; ok {
		channel <- MqttMessage{Data: msg.Payload(), Topic: topic}
	} else {
		// # ile biten subscriptionlar için
		for k, channel := range subscribedChannels {
			if k[len(k)-1:] == "#" && strings.Contains(topic, k[:len(k)-1]) {
				channel <- MqttMessage{Data: msg.Payload(), Topic: topic}
			}
		}
	}
}

func subscribe(conn mqtt.Client, topic string, channel chan MqttMessage) error {
	if token := conn.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return token.Error()
	}
	fmt.Printf("Subscribed to %s\n", topic)
	subscribedChannels[topic] = channel
	return nil
}

func connectMqtt() mqtt.Client {
	host, err := decodeEnv("MQTT_HOST")
	failOnError(err, "set MQTT_HOST env")
	userName, err := decodeEnv("MQTT_USER")
	failOnError(err, "set MQTT_USER env")
	clientID, err := decodeEnv("MQTT_CLIENT_ID")
	failOnError(err, "set MQTT_CLIENT_ID env")
	password, err := decodeEnv("MQTT_PASS")
	failOnError(err, "set MQTT_PASS env")

	opts := mqtt.NewClientOptions().AddBroker(host).SetClientID(clientID)
	opts.Username = userName
	opts.Password = password
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(onMessageHandler)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetOnConnectHandler(onConnectHandler)
	opts.SetCleanSession(true)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return c
}

func decodeEnv(key string) (string, error) {
	result, ok := os.LookupEnv(key)

	if !ok {
		s := fmt.Sprintf("env var not found %s", key)
		return "", errors.New(s)
	}

	return result, nil
}
