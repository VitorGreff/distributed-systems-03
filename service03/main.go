package main

import (
	"log"
	"trab02/rabbitMQ"
)

func main() {
	log.Printf("[ AUTHENTICATION SERVICE N2 ]")
	log.Printf("[ Waiting for messages. Press CTRL+C to terminate the service ]\n\n")
	go rabbitMQ.ReceiveAndGenerateToken()
	go rabbitMQ.ReceiveAndValidateToken()
	select {}
}
