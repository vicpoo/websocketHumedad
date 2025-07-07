//routes.go
package infrastructure

import (
	"github.com/gin-gonic/gin"
)

func SetupHumidityRoutes(r *gin.Engine, hub *Hub) {
	// Ruta WebSocket específica para los datos de humedad
	r.GET("/ws/humidity", hub.HandleWebSocket)

}
