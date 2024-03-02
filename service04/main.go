package main

import (
	"trab02/rabbitMQ"
)

func main() {
	go rabbitMQ.SendProductRequest(4)
	go rabbitMQ.ReadAndSendProduct()
	select {}
}
