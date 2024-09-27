package main

import (
	topic_exchange "github.com/barisaydogdu/MessageQueuesRabbitMQ/internal/messaging"
)

func main() {
	topic_exchange.SendMessage()
}
