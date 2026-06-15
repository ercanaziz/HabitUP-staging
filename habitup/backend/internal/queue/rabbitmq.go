package queue

import (
	"context"
	"encoding/json"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

var conn *amqp.Connection
var ch *amqp.Channel

const HabitCheckedQueue = "habit.checked"

type HabitCheckedEvent struct {
	UserID  string `json:"userId"`
	HabitID string `json:"habitId"`
}

func Connect() {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		url = "amqp://guest:guest@localhost:5672/"
	}

	var err error
	conn, err = amqp.Dial(url)
	if err != nil {
		log.Fatal("RabbitMQ bağlantı hatası:", err)
	}

	ch, err = conn.Channel()
	if err != nil {
		log.Fatal("RabbitMQ kanal hatası:", err)
	}

	_, err = ch.QueueDeclare(HabitCheckedQueue, true, false, false, false, nil)
	if err != nil {
		log.Fatal("RabbitMQ kuyruk hatası:", err)
	}

	log.Println("RabbitMQ bağlantısı başarılı")
}

func PublishHabitChecked(ctx context.Context, userID, habitID string) error {
	body, err := json.Marshal(HabitCheckedEvent{UserID: userID, HabitID: habitID})
	if err != nil {
		return err
	}

	return ch.PublishWithContext(ctx, "", HabitCheckedQueue, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

// StartConsumer alışkanlık tamamlama eventlerini dinler ve loglar.
func StartConsumer() {
	msgs, err := ch.Consume(HabitCheckedQueue, "", true, false, false, false, nil)
	if err != nil {
		log.Println("RabbitMQ consumer hatası:", err)
		return
	}

	go func() {
		for d := range msgs {
			var event HabitCheckedEvent
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Println("Event parse hatası:", err)
				continue
			}
			log.Printf("[RabbitMQ] Alışkanlık tamamlandı — kullanıcı: %s, alışkanlık: %s", event.UserID, event.HabitID)
		}
	}()
}
