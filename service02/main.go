package main

import (
	"log"
	"trab02/rabbitMQ"
)

func main() {
	log.Printf("[ Authentification service has started ]")
	log.Printf("[ Waiting for messages. Press CTRL+C to terminate the service ]\n\n")
	go rabbitMQ.ReceiveAndGenerateToken()
	rabbitMQ.ReceiveAndValidateToken()
}
