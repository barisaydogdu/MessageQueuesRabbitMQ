package direct

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

//	func failOnErrorReceive(err error, msg string) {
//		if err != nil {
//			log.Panicf("%s:%s", msg, err)
//		}
//	}
func receiveMessage() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs_direct", //Oluşturulan exchange adı
		"direct",      //type
		true,          //durable- exchangenin kalıcılık. True değerinde rabbitMQ sunucusu yeniden başladığında exchange kaybolmaz
		false,         //auto deleted- true değerinde herhangi bir kuyruğa bağlı olmadığında silinir
		false,         //internal- içsel olup olmadığını belirtir. false ise exchange dışarıdan gelen mesajlar için de kullanılabilir
		false,         //istemcinin exchanenin oluştutrulması tamamlanana kadar beklemek isteyip istemediğini belirtir. false olduğunda exhange başarılı bir şekilde oluşturulup oluşturuşmadığıan dair bir onay bekler
		nil,           //arguments
	)
	failOnError(err, "Failed to declare an exchange")
	//5 saniye içinde tamamlanmazsa iptal et
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := bodyFrom(os.Args)
	err = ch.PublishWithContext(ctx,
		"logs_direct",         //exchange
		severityFrom(os.Args), // Routing Keys
		false,                 // mandatory -herhangi bir kuyruk bu routing key ile eşleşmiyorsa hata alıp almayacağını belirler. False olduğunda hata fırlatmaz
		false,                 //immediate- Mesajın hemen tüketilip tüketilmemesi gerektiğini belirler. False olduğunda hemen tüketilmeyebilir. Önce bir kuyrukta saklabilir.
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf("[x] Sent %s", body)
}

func bodyFrom(args []string) string {
	var s string
	if len(args) < 3 || os.Args[2] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[2:], " ")
	}
	return s
}

func severityFrom(args []string) string {

	var s string
	if len(args) < 2 || os.Args[1] == "" {
		s = "info"
	} else {
		s = os.Args[1]
	}
	return s
}
