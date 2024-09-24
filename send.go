package main

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s:%s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	//Create a channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", //name
		false,   //sürdürülebilir
		false,   //kullanıldığında sil
		false,   //özel
		false,   //bekleme yok
		nil,     //argümanlar
	)
	failOnError(err, "Failed to declare a queue")

	//işlem 5 saniye içinde tamamlanmazsa iptal et
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Hello World"
	err = ch.PublishWithContext(ctx,
		"",     //exhange
		q.Name, //yönlendirme anahtarı
		false,  //zorunlu
		false,  //acil
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf("[x] Sent %s\n", body)
}
