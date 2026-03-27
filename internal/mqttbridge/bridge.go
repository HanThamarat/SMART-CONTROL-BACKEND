package mqttbridge

import (
	"log"
	"os"
	"strings"

	"github.com/HanThamarat/SMART-CONTROL-BACKEND/internal/socket"
	mqttcon "github.com/HanThamarat/SMART-CONTROL-BACKEND/pkg/mqttCon"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func Setup(socketServer *socket.Server) mqtt.Client {
	mqttClient := mqttcon.MqttConnection()
	socketServer.SetMQTTPublisher(mqttcon.NewPublisher(mqttClient))
	subscribeTopics(mqttClient, socketServer)

	return mqttClient
}

func subscribeTopics(mqttClient mqtt.Client, socketServer *socket.Server) {
	for _, topic := range topicsFromEnv() {
		topic := topic
		tokenSub := mqttClient.Subscribe(topic, 0, func(c mqtt.Client, m mqtt.Message) {
			log.Printf("mqtt message topic=%s payload=%s", m.Topic(), string(m.Payload()))
			socketServer.BroadcastMQTTMessage(m.Topic(), m.Payload())
		})
		tokenSub.Wait()

		if err := tokenSub.Error(); err != nil {
			log.Fatalf("Failed to subscribe topic %s: %v", topic, err)
		}

		log.Printf("mqtt subscribed topic=%s", topic)
	}
}

func topicsFromEnv() []string {
	rawTopics := strings.TrimSpace(os.Getenv("MQTT_SUBSCRIBE_TOPICS"))
	if rawTopics == "" {
		rawTopics = "TEST/MQTT"
	}

	parts := strings.Split(rawTopics, ",")
	topics := make([]string, 0, len(parts))
	for _, part := range parts {
		topic := strings.TrimSpace(part)
		if topic == "" {
			continue
		}

		topics = append(topics, topic)
	}

	return topics
}
