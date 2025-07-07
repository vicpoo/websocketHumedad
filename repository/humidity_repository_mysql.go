// humidity_repository_mysql.go
package repository

import (
	"database/sql"
	"fmt"

	"github.com/vicpoo/websocketHumedad/Humedad/domain"
	"github.com/vicpoo/websocketHumedad/Humedad/domain/entities"
	"github.com/vicpoo/websocketHumedad/core"
)

type humidityRepositoryMySQL struct {
	db *sql.DB
}

func NewHumidityRepositoryMySQL() domain.HumidityRepository {
	return &humidityRepositoryMySQL{
		db: core.GetBD(),
	}
}

func (r *humidityRepositoryMySQL) Save(data entities.HumidityTemperatureData) error {
	var sensorID int
	err := r.db.QueryRow("SELECT id FROM sensors WHERE name = ?", data.Sensor).Scan(&sensorID)
	if err != nil {
		return fmt.Errorf("no se encontró el sensor '%s': %v", data.Sensor, err)
	}

	_, err = r.db.Exec(`
		INSERT INTO sensor_readings (
			sensor_id, temperature, humidity, recorded_at
		) VALUES (?, ?, ?, FROM_UNIXTIME(?))`,
		sensorID, data.Temperature, data.Humidity, data.Timestamp)

	if err != nil {
		return fmt.Errorf("error al insertar en sensor_readings: %v", err)
	}

	return nil
}

func (r *humidityRepositoryMySQL) GetAll() ([]entities.HumidityTemperatureData, error) {
	rows, err := r.db.Query(`
		SELECT s.name, sr.temperature, sr.humidity, UNIX_TIMESTAMP(sr.recorded_at), s.location
		FROM sensor_readings sr
		JOIN sensors s ON sr.sensor_id = s.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []entities.HumidityTemperatureData
	for rows.Next() {
		var d entities.HumidityTemperatureData
		var timestamp int64

		err := rows.Scan(&d.Sensor, &d.Temperature, &d.Humidity, &timestamp, &d.Location)
		if err != nil {
			return nil, err
		}
		d.Timestamp = timestamp
		result = append(result, d)
	}

	return result, nil
}
