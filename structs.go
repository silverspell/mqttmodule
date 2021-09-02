package mqttmodule

type MqttMessage struct {
	Topic string
	Data  []byte
}
