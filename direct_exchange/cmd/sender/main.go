package main

import (
	direct_exchange "github.com/barisaydogdu/MessageQueuesRabbitMQ/internal/messaging"
)

func main() {
	direct_exchange.SendMessage()
}

//cmd içerisinde iki ayrı app kullanmak doğru mu?
//Github repolarında böyle kullanıldıklarını gördüm. Bunu görmeden öncesinde waitgroup ile
//yapmayı düşünüyordum
