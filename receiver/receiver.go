package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {

	if err != nil {
		log.Panicf("%s:%s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", //name
		false,   //sunucu yeniden başladığında kaybolmaması için kalıcılığı temsil eder
		false,   // kuyruk kullanılmadığında otomatik olarak silinsin mi ?
		false,   //exclusive- özel - false ise birden fazla bağlantı bu kuyruğu kullabilir. Farklı uygulamalar ve iş parçacıkları bu kuyruğa mesaj gönderebilir veya alabilir
		false,   //no wait- false ise RabbitMQ sunucunun kuyruğun başarıyla oluşup oluşturulmadığını belirlemke için bir yanıt göndeirr
		nil,     // argümanlar
	)
	failOnError(err, "Failed to declare a queue")

	//Bir goroutine içindeki bir kanaldan okuyacağız
	//tüketici
	msgs, err := ch.Consume(
		q.Name, //queue
		"",     //consumer Mesajların hangi tüketici tarafından tüketileceğini belirtir
		true,   // auto-ack true ile birlikte mesajları otomatik onaylar. Değer false olursa işlendikten sonra manuel olarak onaylanmalıdır.
		false,  //exclusive- özel false ile birlikte kuyruktan birden fazla bağlantı ile mesaj alınabilir
		false,  //no local- false ise tüketici aynı bağlantı üzerinden mesaj gönderiyorsa bu mesajı alabilir
		false,  //no-wait RabbitMQ sunucunun yanıt verip vermeyeceğini belirtir.
		nil,    //Argümanlar
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	log.Printf("[*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
