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
		log.Println("RABBITMQ_URL ayarlanmamış — RabbitMQ devre dışı")
		return
	}

	var err error
	conn, err = amqp.Dial(url)
	if err != nil {
		log.Printf("RabbitMQ bağlantı hatası (devam ediliyor): %v", err)
		return
	}

	ch, err = conn.Channel()
	if err != nil {
		log.Printf("RabbitMQ kanal hatası (devam ediliyor): %v", err)
		return
	}

	_, err = ch.QueueDeclare(HabitCheckedQueue, true, false, false, false, nil)
	if err != nil {
		log.Printf("RabbitMQ kuyruk hatası (devam ediliyor): %v", err)
		return
	}

	log.Println("RabbitMQ bağlantısı başarılı")
}


func PublishHabitChecked(ctx context.Context, userID, habitID string) error {
	if ch == nil {
		log.Println("RabbitMQ kanalı yok — event atlanıyor")
		return nil
	}
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
	if ch == nil {
		log.Println("RabbitMQ kanalı yok — consumer başlatılmıyor")
		return
	}
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
