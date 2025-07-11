//main.go
package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/vicpoo/websocketHumedad/core"
	"github.com/vicpoo/websocketHumedad/Humedad/infrastructure"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	// ✅ Inicializa la base de datos
	core.InitDB()

	r := gin.Default()

	// CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		log.Printf("Request: %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
	})

	// WebSocket hub
	hub := infrastructure.NewHub()
	go hub.Run()

	// Servicio de mensajería (RabbitMQ)
	messagingService := infrastructure.NewMessagingService(hub)
	defer messagingService.Close()

	// Rutas WebSocket para humedad
	infrastructure.SetupHumidityRoutes(r, hub)

	// Inicia consumidor RabbitMQ para sensor de humedad
	if err := messagingService.ConsumeHumidityMessages(); err != nil {
		log.Fatalf("Failed to start Humidity consumer: %v", err)
	}

	// Apagado controlado
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Inicia servidor
	go func() {
		if err := r.Run(":8001"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Println("Server started on port 8001")
	log.Println("Humidity RabbitMQ consumer started")

	<-sigChan
	log.Println("Shutting down server...")
}
