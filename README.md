# Mqtt Module

## Gereksinimler
- go get github.com/silverspell/mqttmodule
- MQTT_HOST, MQTT_USER, MQTT_PASS, MQTT_CLIENT_ID env değişkenleri

## Subscribe işlemleri 

Mqtt subscribe ederken her bir subscription için bir MqttMessage tipinde channel tanımlamalısınız. Topic ve channelları aynı sıra ile Subscribe metoduna vermeniz gerekmekte. Teknik olarak bir uygulamada en fazla 2 ya da 3 subscription olacağı düşünülürse (wildcards bro) makul bir çözüm.

```go
func subscribe() {
    topic := "devices/lamps/#"
	listenChan := make(chan mqttmodule.MqttMessage)
	chans := make([]chan mqttmodule.MqttMessage, 1)
	chans[0] = listenChan

	go mqttmodule.Subscribe([]string{topic}, chans)

	for {
		msg := <-listenChan
		fmt.Printf("%s %s\n", msg.Topic, string(msg.Data))
	}
}

```

## Publish

```go
func publish() {
    topic := "devices/lamps/livingroom"
	listenChan := make(chan mqttmodule.MqttMessage)
	

	go mqttmodule.Publish(listenChan)

	for i := 0; i< 10; i++{ 
		msg := mqttmodule.MqttMessage{Topic: topic, Data: []byte("Hello")}
		listenChan <- msg
	}
}

```
