//messaging_service.go
package infrastructure

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
	"github.com/vicpoo/websocketHumedad/Humedad/application"
	"github.com/vicpoo/websocketHumedad/Humedad/domain/entities"
	"github.com/vicpoo/websocketHumedad/repository"
)

type MessagingService struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	hub  *Hub
}

func NewMessagingService(hub *Hub) *MessagingService {
	conn, err := amqp.Dial("amqp://reyhades:reyhades@44.219.123.4:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
		return nil
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
		return nil
	}

	err = ch.ExchangeDeclare(
		"amq.topic",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare an exchange: %s", err)
		return nil
	}

	return &MessagingService{
		conn: conn,
		ch:   ch,
		hub:  hub,
	}
}

func (ms *MessagingService) ConsumeHumidityMessages() error {
	q, err := ms.ch.QueueDeclare(
		"sensor_data",
		true, false, false, false, nil,
	)
	if err != nil {
		return err
	}

	err = ms.ch.QueueBind(
		q.Name,
		"sensor_inte",
		"amq.topic",
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := ms.ch.Consume(
		q.Name, "", false, false, false, false, nil,
	)
	if err != nil {
		return err
	}

	// Inicializa el repositorio y use case
	repo := repository.NewHumidityRepositoryMySQL()
	useCase := application.NewHumidityUseCase(repo)

	go func() {
		for msg := range msgs {
			log.Printf("Mensaje sensor_inte recibido: %s", string(msg.Body))

			var payload struct {
				Sensor            string  `json:"sensor"`
				Temperatura       float64 `json:"temperatura"`
				Humedad           float64 `json:"humedad"`
				UnidadTemperatura string  `json:"unidad_temperatura"`
				UnidadHumedad     string  `json:"unidad_humedad"`
				Timestamp         int64   `json:"timestamp"`
				Ubicacion         string  `json:"ubicacion"`
			}

			if err := json.Unmarshal(msg.Body, &payload); err != nil {
				log.Printf("Error al parsear mensaje JSON: %v", err)
				msg.Nack(false, false)
				continue
			}

			data := entities.NewHumidityTemperatureData(
				payload.Sensor,
				payload.Temperatura,
				payload.Humedad,
				payload.UnidadTemperatura,
				payload.UnidadHumedad,
				payload.Timestamp,
				payload.Ubicacion,
			)

			if err := useCase.SaveHumidityData(*data); err != nil {
				log.Printf("Error al guardar en BD: %v", err)
			} else {
				log.Printf("Datos guardados correctamente")
			}

			// Enviar a WebSocket
			ms.hub.broadcast <- msg.Body
			msg.Ack(false)
		}
	}()

	return nil
}

func (ms *MessagingService) Close() {
	if ms.ch != nil {
		ms.ch.Close()
	}
	if ms.conn != nil {
		ms.conn.Close()
	}
}
