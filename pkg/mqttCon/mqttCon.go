package mqttcon

import (
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

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
