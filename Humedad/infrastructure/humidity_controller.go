//humidity_controller.go
package infrastructure

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/websocketHumedad/Humedad/application"
)

type HumidityController struct {
	useCase *application.HumidityUseCase
}

func NewHumidityController(useCase *application.HumidityUseCase) *HumidityController {
	return &HumidityController{useCase: useCase}
}

func (hc *HumidityController) GetAll(c *gin.Context) {
	data, err := hc.useCase.GetAllHumidityData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
