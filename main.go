package main

import (
	// "fmt"
	"log"
	"os"

	"github.com/HanThamarat/SMART-CONTROL-BACKEND/internal/domain"
	"github.com/HanThamarat/SMART-CONTROL-BACKEND/internal/handler"
	"github.com/HanThamarat/SMART-CONTROL-BACKEND/internal/repositories"
	"github.com/HanThamarat/SMART-CONTROL-BACKEND/internal/router"
	"github.com/HanThamarat/SMART-CONTROL-BACKEND/internal/usecase"
	"github.com/HanThamarat/SMART-CONTROL-BACKEND/pkg/database"
	initial "github.com/HanThamarat/SMART-CONTROL-BACKEND/pkg/initialize"
	loadend "github.com/HanThamarat/SMART-CONTROL-BACKEND/pkg/loadEnd"

	// mqttcon "github.com/HanThamarat/SMART-CONTROL-BACKEND/pkg/mqttCon"
	// mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)


func main() {
	loadend.LoadEnv();

	db := database.InitDB();

	err := db.AutoMigrate(
		&domain.User{},
	);

	if err != nil {
		panic("Could not migrate database: " + err.Error());
	}

	initial.UserInit(db);

	// mqttClient := mqttcon.MqttConnection();

	// topic := "smart/control"

	// text := "Hello from Go!"
	// token := mqttClient.Publish("smart/control", 0, false, text);
	// token.Wait();

	// tokenSub := mqttClient.Subscribe(topic, 0, func(c mqtt.Client, m mqtt.Message) {
	// 	fmt.Printf("✅ Received message: %s from topic: %s\n", m.Payload(), m.Topic());
	// });
	// tokenSub.Wait();

	authRepo 	:= repositories.NewGormAuthRepository(db);
	authUc 		:= usecase.NewAuthUsecase(authRepo);
	authHdl 	:= handler.NewAuthHandler(authUc);

	app := fiber.New();
	app.Use(logger.New());
	app.Use(cors.New(cors.Config{
		AllowOrigins:  "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, HEAD, PUT, DELETE, PATCH",
		AllowCredentials: false,
	}));

	router.SetupRoutes(
		app,
		authHdl,
	);

	log.Fatal(app.Listen(os.Getenv("Port")));
}