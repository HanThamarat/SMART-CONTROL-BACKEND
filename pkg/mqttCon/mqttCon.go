package mqttcon

import (
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Publisher struct {
	client mqtt.Client
}

func MqttConnection() mqtt.Client {

	// setting client option connection to mqtt broker.
	opts := mqtt.NewClientOptions()
	opts.AddBroker(os.Getenv("MQTT_BROKER"))
	opts.SetClientID(os.Getenv("MQTT_CLIENT_ID"))
	opts.SetUsername(os.Getenv("MQTT_USERNAME"))
	opts.SetPassword(os.Getenv("MQTT_PASSWORD"))
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(1 * time.Minute)
	opts.SetCleanSession(false) // Remembers subscriptions on reconnect
	opts.SetOnConnectHandler(func(c mqtt.Client) {
		fmt.Println("✨ Connection established/re-established")
	})

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect: %v", token.Error())
	}

	fmt.Println("✅ MQTT Broker connected.")

	return client
}

func NewPublisher(client mqtt.Client) *Publisher {
	return &Publisher{client: client}
}

func (p *Publisher) Publish(topic string, qos byte, retained bool, payload interface{}) error {
	if p == nil || p.client == nil {
		return fmt.Errorf("mqtt client is not configured")
	}

	token := p.client.Publish(topic, qos, retained, payload)
	token.Wait()

	return token.Error()
}
