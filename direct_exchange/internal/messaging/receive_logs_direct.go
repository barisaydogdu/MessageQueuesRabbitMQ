package direct_exchange

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s:%s", msg, err)
	}
}

func ReceiveMessage() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs_direct", //nane
		"direct",      //exchange type
		true,          //durable- yeniden başlatıldığında exhange kaybolmaz
		false,         //auto-deleted
		false,         //internal - içsel
		false,         //nowait
		nil,
	)
	failOnError(err, "Failed to declare an exchange")
	q, err := ch.QueueDeclare(
		"",    //name
		false, //durable
		false, //delete when unused
		true,  //exclusive- true olduğunda kuyruk sadece bu kanal ile kullanılabilir.
		false, //no-wait
		nil,   //arguments
	)
	failOnError(err, "Failed to declare a queue")
	if len(os.Args) < 2 {
		log.Printf("Usage:%s [info] [warning] [error]", os.Args[0])
		os.Exit(0)
	}
	for _, s := range os.Args[1:] {
		log.Printf("Binding queue %s to exchange %s with routing key %s", q.Name, "logs_direct", s)
		err = ch.QueueBind(
			q.Name,        //Queue name
			s,             //Routing key
			"logs_direct", //exchange
			false,         //no-wait
			nil)           //args)
		failOnError(err, "Failed to bind a queue")
	}
	//tüketici
	msgs, err := ch.Consume(
		q.Name, //queue
		"",
		true,  // auto-ack- Mesajı otomatik olarak onaylamak
		false, //exclusive- sadece bu tüketici mi dinleyecek
		false, //nolocal- bu bağlantı üzerinden gönderilen messjların bu bağlatıdan dinleyen bir tüketiciye iletilip iletilmediğini kontrol eder.
		false, //nowait- RabbitMQ'dan bir onay mesajı beklenir ve bağlantı kurulduğu zaman bir yanıt alınır.
		nil,
	)
	failOnError(err, "Failed to register a consumer")
	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("[x] %s", d.Body)
		}
	}()
	log.Printf("[*] Waiting for logs. To Exit press CTRL+C")
	<-forever

}
