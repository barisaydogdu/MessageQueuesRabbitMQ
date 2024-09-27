package main

import (
	fanout_exchange "github.com/barisaydogdu/MessageQueuesRabbitMQ/internal/messaging"
)

func main() {

	fanout_exchange.SendMessage()
}
