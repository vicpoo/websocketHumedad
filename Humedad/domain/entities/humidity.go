//humidity.go
package entities

import "time"

type HumidityTemperatureData struct {
	ID                int       `json:"id"`
	Sensor            string    `json:"sensor"`
	Temperature       float64   `json:"temperatura"`
	Humidity          float64   `json:"humedad"`
	TemperatureUnit   string    `json:"unidad_temperatura"`
	HumidityUnit      string    `json:"unidad_humedad"`
	Timestamp         int64     `json:"timestamp"`
	Location          string    `json:"ubicacion"`
	CreatedAt         time.Time `json:"created_at"`
}

func NewHumidityTemperatureData(
	sensor string,
	temperature float64,
	humidity float64,
	tempUnit string,
	humidUnit string,
	timestamp int64,
	location string,
) *HumidityTemperatureData {
	return &HumidityTemperatureData{
		Sensor:           sensor,
		Temperature:      temperature,
		Humidity:         humidity,
		TemperatureUnit:  tempUnit,
		HumidityUnit:     humidUnit,
		Timestamp:        timestamp,
		Location:         location,
		CreatedAt:        time.Now(),
	}
}
