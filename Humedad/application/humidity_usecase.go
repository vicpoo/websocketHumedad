//humidity_usecase.go
package application

import (
	"github.com/vicpoo/websocketHumedad/Humedad/domain"
	"github.com/vicpoo/websocketHumedad/Humedad/domain/entities"
)

type HumidityUseCase struct {
	repo domain.HumidityRepository
}

func NewHumidityUseCase(repo domain.HumidityRepository) *HumidityUseCase {
	return &HumidityUseCase{repo: repo}
}

func (uc *HumidityUseCase) SaveHumidityData(data entities.HumidityTemperatureData) error {
	return uc.repo.Save(data)
}

func (uc *HumidityUseCase) GetAllHumidityData() ([]entities.HumidityTemperatureData, error) {
	return uc.repo.GetAll()
}
