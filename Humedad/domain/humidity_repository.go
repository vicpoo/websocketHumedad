package domain

import "github.com/vicpoo/websocketHumedad/Humedad/domain/entities"

type HumidityRepository interface {
	Save(data entities.HumidityTemperatureData) error
	GetAll() ([]entities.HumidityTemperatureData, error)
}
