package main

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s:%s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs_topic", //name
		"topic",      //type
		true,         //durable
		false,        //auto deleted
		false,        //internal
		false,        //nowait
		nil,
	)
	failOnError(err, "Failed to declare an exhange")

	q, err := ch.QueueDeclare(
		"",    //name
		false, //durable
		false, //delete when unused
		true,  //exclusive kuyruk sadece bu bağlantıdan erişebilir.
		false, //no wait
		nil,   //arguments
	)
	failOnError(err, "Failed to declare a queue")
	if len(os.Args) < 2 {
		log.Printf("Usage:%s [binding_key]....", os.Args[0])
		os.Exit(0)
	}

	for _, s := range os.Args[1:] {
		log.Printf("Binding queue %s to exchange %s with routing key %s",
			q.Name, "logs_topic", s)

		err = ch.QueueBind(
			q.Name,       //queue name
			s,            //routing key
			"logs_topic", //exchane
			false,
			nil)
		failOnError(err, "Failed to bind a queue")
	}
	msgs, err := ch.Consume(
		q.Name,
		"",
		true,  //autoACK -mesjaların otomatik onaylanıp onaylanmadığını belirler. true hemen onaylar
		false, //exclusive -Birden fazla bağlantı ile mesaj alabilir
		false, //nolocal	true ise tüketici kendi gönderdiği mesajları almaz.
		false, //no wait //sunucu tüketici oluşturma isteğine yanıt verecektir.
		nil,   //args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("[x] %s", d.Body)
		}
	}()
	log.Printf("[*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
