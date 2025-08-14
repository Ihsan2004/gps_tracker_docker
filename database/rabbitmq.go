package database

import (
	"GpsTracker2/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"
)

const (
	rabbitMQURL = "amqp://guest:guest@localhost:5672/"
	queueName   = "location_history"
)

var (
	rabbitConn *amqp.Connection
	rabbitChan *amqp.Channel
)

// initRabbitMQ initializes the RabbitMQ connection and channel
func initRabbitMQ() error {
	var err error
	rabbitConn, err = amqp.Dial(rabbitMQURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	rabbitChan, err = rabbitConn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}

	_, err = rabbitChan.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	return nil
}

// PublishMessage publishes a message to RabbitMQ
func PublishMessage(body []byte) error {
	if rabbitChan == nil {
		if err := initRabbitMQ(); err != nil {
			return err
		}
	}

	err := rabbitChan.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		// Try to reinitialize connection and retry once
		if err := initRabbitMQ(); err != nil {
			return fmt.Errorf("failed to reinitialize RabbitMQ: %w", err)
		}

		err = rabbitChan.Publish(
			"",        // exchange
			queueName, // routing key
			false,     // mandatory
			false,     // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			},
		)
		if err != nil {
			return fmt.Errorf("failed to publish message after reconnection: %w", err)
		}
	}

	log.Printf("Successfully published message to RabbitMQ")
	return nil
}

// processMessage handles consuming messages from RabbitMQ
func processMessage() {
	if rabbitChan == nil {
		if err := initRabbitMQ(); err != nil {
			log.Fatalf("Failed to initialize RabbitMQ: %v", err)
		}
	}

	msgs, err := rabbitChan.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			var location models.Location
			if err := json.Unmarshal(msg.Body, &location); err != nil {
				log.Printf("Error decoding message: %v", err)
				continue
			}

			// Save to MongoDB
			_, err = LocationCollection.InsertOne(context.TODO(), bson.M{
				"device_id": location.DeviceID,
				"longitude": location.Longitude,
				"latitude":  location.Latitude,
				"timestamp": time.Now(), // Adding timestamp for when the location was saved
			})
			if err != nil {
				log.Printf("Error saving to MongoDB: %v", err)
			} else {
				log.Printf("Location saved to MongoDB for device %d", location.DeviceID)
			}
		}
	}()

	log.Printf("Location consumer started. Waiting for messages...")
	<-forever
}

// StartConsumer initializes and starts the RabbitMQ consumer
func StartConsumer() {
	ConnectMongo()
	processMessage()
}

// CloseConnections closes RabbitMQ connections
func CloseConnections() {
	if rabbitChan != nil {
		rabbitChan.Close()
	}
	if rabbitConn != nil {
		rabbitConn.Close()
	}
}
