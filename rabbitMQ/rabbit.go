package rabbitMQ

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"
	"trab02/models"
	"trab02/service01/database"
	products_repository "trab02/service03/repositories"
	tokenPkg "trab02/token"

	amqp "github.com/rabbitmq/amqp091-go"
)

func prepareRabbitMQ() (*amqp.Connection, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, errors.New("failed to connect to RabbitMQ: " + err.Error())
	}
	return conn, nil
}

func SendAndConsumeToken(token string, userID uint64) error {
	conn, err := prepareRabbitMQ()
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return errors.New("failed to open a channel: " + err.Error())
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"token_validation",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.New("failed to declare a queue: " + err.Error())
	}

	message := token + "," + strconv.FormatUint(userID, 10)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = ch.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		return errors.New("failed to publish a message: " + err.Error())
	}

	q2, err := ch.QueueDeclare(
		"result_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.New("failed to declare a queue: " + err.Error())
	}

	msgs, err := ch.Consume(
		q2.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.New("failed to register a consumer: " + err.Error())
	}

	var response string
	for d := range msgs {
		log.Println("Response", string(d.Body))
		response = string(d.Body)
		break
	}
	if response == "invalid" {
		return errors.New("token inválido")
	}
	return nil
}

func ReceiveAndValidateToken() {
	conn, err := prepareRabbitMQ()
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Println(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"token_validation",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			parts := strings.Split(string(d.Body), ",")
			if len(parts) != 2 {
				log.Println("Invalid message received")
				sendValidationResponse(ch, errors.New("invalid message"))
				continue
			}

			token := parts[0]
			userID, err := strconv.ParseUint(parts[1], 10, 64)
			if err != nil {
				log.Println("Invalid userID received")
				sendValidationResponse(ch, errors.New("invalid userID"))
				continue
			}

			log.Printf("Received a message: %s", d.Body)
			result := tokenPkg.ValidateToken(token, userID)
			sendValidationResponse(ch, result)
		}
	}()
	<-forever
}

func sendValidationResponse(ch *amqp.Channel, result error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response := "valid"
	if result != nil {
		response = "invalid"
	}

	err := ch.PublishWithContext(
		ctx,
		"",
		"result_queue",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(response),
		},
	)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Response sent: %v\n\n", result)
}

func SendTokenGenerationRequest(userID uint64) (string, error) {
	conn, err := prepareRabbitMQ()
	if err != nil {
		return "", err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return "", err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"token_generation",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = ch.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(strconv.FormatUint(userID, 10)),
		},
	)
	if err != nil {
		return "", err
	}

	q2, err := ch.QueueDeclare(
		"token_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return "", err
	}

	msgs, err := ch.Consume(
		q2.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return "", err
	}

	var token string
	for d := range msgs {
		token = string(d.Body)
		break
	}
	if token != "invalid" {
		return token, nil
	}
	return "", errors.New("token inválido")
}

func ReceiveAndGenerateToken() {
	conn, err := prepareRabbitMQ()
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Println(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"token_generation",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			userID, err := strconv.ParseUint(string(d.Body), 10, 64)
			if err != nil {
				log.Println("Invalid userID received")
				sendToken(ch, "token_queue", "invalid")
				continue
			}
			log.Printf("Received a message from id: %s", d.Body)
			token, err := tokenPkg.GenerateToken(userID)
			if err != nil {
				sendToken(ch, "token_queue", "invalid")
				continue
			}
			sendToken(ch, "token_queue", token)
		}
	}()
	<-forever
}

func sendToken(ch *amqp.Channel, queueName string, token string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := ch.PublishWithContext(
		ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(token),
		},
	)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Token sent: %s\n\n", token)
}

// 3 service
func SendProductRequest(productId uint64) (*models.Product, error) {
	conn, err := prepareRabbitMQ()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return nil, errors.New("failed to open a channel: " + err.Error())
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"product_request",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, errors.New("failed to declare a queue: " + err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = ch.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(strconv.FormatUint(productId, 10)),
		},
	)
	if err != nil {
		return nil, errors.New("failed to publish a message: " + err.Error())
	}

	q2, err := ch.QueueDeclare(
		"product_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, errors.New("failed to declare a queue: " + err.Error())
	}

	msgs, err := ch.Consume(
		q2.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, errors.New("failed to register a consumer: " + err.Error())
	}

	var product models.Product
	for d := range msgs {
		log.Println("Response", string(d.Body))
		if err := json.Unmarshal(d.Body, &product); err != nil {
			return nil, err
		} else {
			break
		}
	}
	return &product, nil
}
func ReadAndSendProduct() {
	conn, err := prepareRabbitMQ()
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Println(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"product_request",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			userID, err := strconv.ParseUint(string(d.Body), 10, 64)
			if err != nil {
				log.Println("Invalid userID received")
				sendProduct(ch, nil, "product_queue")
				continue
			}
			log.Printf("Received a message from id: %s", d.Body)
			db, err := database.InitMySqlConn()
			if err != nil {
				sendProduct(ch, nil, "product_queue")

			}
			repo := products_repository.NewProductRepository(db)
			product, err := repo.GetProduct(userID)
			repo.DecrementProductQuantity(userID)
			if err != nil {
				sendProduct(ch, nil, "product_queue")
				continue
			}
			sendProduct(ch, &product, "product_queue")
		}
	}()
	<-forever
}

func sendProduct(ch *amqp.Channel, product *models.Product, queueName string) {
	log.Println(product)
	log.Println(*product)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	body, err := json.Marshal(*product)
	if err != nil {
		log.Println(err)
	}
	err = ch.PublishWithContext(
		ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
	if err != nil {
		log.Println(err)
	}
	// log.Println("Product sent:", product)
}
