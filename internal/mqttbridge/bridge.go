package mqttbridge

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/HanThamarat/SMART-CONTROL-BACKEND/internal/domain"
	"github.com/HanThamarat/SMART-CONTROL-BACKEND/internal/socket"
	"github.com/HanThamarat/SMART-CONTROL-BACKEND/internal/types"
	mqttcon "github.com/HanThamarat/SMART-CONTROL-BACKEND/pkg/mqttCon"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gorm.io/gorm"
)

func Setup(socketServer *socket.Server, db *gorm.DB) mqtt.Client {
	mqttClient := mqttcon.MqttConnection()
	socketServer.SetMQTTPublisher(mqttcon.NewPublisher(mqttClient))
	subscribeTopics(mqttClient, socketServer, db)

	return mqttClient
}

func subscribeTopics(mqttClient mqtt.Client, socketServer *socket.Server, db *gorm.DB) {
	go func () {
		subscribed := make(map[string]bool)

		for {
			topics := topicsFromEnv(db)

			for _, topic := range topics {
				if subscribed[topic] {
					continue
				}

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
			
			current := make(map[string]bool)
			for _, t := range topics {
				current[t] = true
			}

			for t := range subscribed {
				if !current[t] {
					token := mqttClient.Unsubscribe(t)
					token.Wait()

					if err := token.Error(); err != nil {
						log.Printf("❌ Failed to unsubscribe topic %s: %v", t, err)
						continue
					}

					delete(subscribed, t)
					log.Printf("🗑️ unsubscribed topic=%s", t)
				}
			}

			time.Sleep(5 * time.Second)
		}
	}()
}

func topicsFromEnv(db *gorm.DB) []string {
	var widgetModel domain.Widget;
	var widgetName	[]types.TopicQuery;

	if err := db.Select("widget_name").Where("deleted_at IS NULL").Model(&widgetModel).Scan(&widgetName).Error; err != nil {
		fmt.Println("getting widget for create topic name failed : ", err);
	}

	topics := make([]string, 0, len(widgetName))
	for _, part := range widgetName {
		topic := strings.TrimSpace(part.WidgetName)
		if topic == "" {
			continue
		}

		topics = append(topics, topic)
	}

	return topics
}
