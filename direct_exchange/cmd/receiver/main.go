package main

import (
	direct_exchange "github.com/barisaydogdu/MessageQueuesRabbitMQ/internal/messaging"
)

func main() {
	direct_exchange.ReceiveMessage()
}
