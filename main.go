package main

import (
	"fmt"

	loadend "github.com/HanThamarat/SMART-CONTROL-BACKEND/pkg/loadEnd"
	mqttcon "github.com/HanThamarat/SMART-CONTROL-BACKEND/pkg/mqttCon"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)


func main() {
	loadend.LoadEnv();

	mqttClient := mqttcon.MqttConnection();

	topic := "smart/control"

	text := "Hello from Go!"
	token := mqttClient.Publish("smart/control", 0, false, text);
	token.Wait();

	tokenSub := mqttClient.Subscribe(topic, 0, func(c mqtt.Client, m mqtt.Message) {
		fmt.Printf("✅ Received message: %s from topic: %s\n", m.Payload(), m.Topic());
	});
	tokenSub.Wait();

	keepAlive := make(chan bool)
    fmt.Println("🚀 Service is running. Press CTRL+C to stop.")
    <-keepAlive
}