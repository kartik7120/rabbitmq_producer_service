package producers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	rabbitmq_producer "github.com/kartik7120/booking_rabbitmq_producer_service/cmd/grpcServer"
	"github.com/kartik7120/booking_rabbitmq_producer_service/cmd/models"
	"github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	Conn *amqp091.Channel
}

func NewProducer(conn *amqp091.Channel) *Producer {
	return &Producer{
		Conn: conn,
	}
}

func (p *Producer) Payment_Service_Producer(payload models.Payment) error {

	// Declare exchange

	err := p.Conn.ExchangeDeclare(
		"payment_success_exchange",
		"direct", // exchange type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)

	if err != nil {
		return err
	}

	q, err := p.Conn.QueueDeclare(
		"payment_service_success",
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		return err
	}

	// Bind the queue to the exchange

	err = p.Conn.QueueBind(
		q.Name,                     // queue name
		"payment_success_key",      // routing key
		"payment_success_exchange", // exchange name
		false,                      // no-wait
		nil,                        // arguments
	)

	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	err = p.Conn.PublishWithContext(
		ctx,
		"payment_success_exchange",
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
			Timestamp:   time.Now(),
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func (p *Producer) Payment_Service_Failure_Producer(payload models.Payment) error {

	// Declare exchange

	err := p.Conn.ExchangeDeclare(
		"payment_failure_exchange",
		"direct", // exchange type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)

	if err != nil {
		return err
	}

	q, err := p.Conn.QueueDeclare(
		"payment_service_failure",
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		return err
	}

	// Bind the queue to the exchange

	err = p.Conn.QueueBind(
		q.Name,                     // queue name
		"payment_failure_key",      // routing key
		"payment_failure_exchange", // exchange name
		false,                      // no-wait
		nil,                        // arguments
	)

	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	err = p.Conn.PublishWithContext(
		ctx,
		"payment_failure_exchange",
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
			Timestamp:   time.Now(),
		},
	)

	if err != nil {
		return err
	}

	return nil
}

// Create producer for lock seats and unlock seats

func (p *Producer) Lock_Seats(seatsIds []int) error {

	err := p.Conn.ExchangeDeclare(
		"lock_seats",
		"direct",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	q, err := p.Conn.QueueDeclare(
		"lock_seats_queue",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	err = p.Conn.QueueBind(
		q.Name,
		"lock_seats_key",
		"lock_seats",
		false,
		nil,
	)

	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	bodyBytes, err := json.Marshal(seatsIds)

	if err != nil {
		return err
	}

	err = p.Conn.PublishWithContext(
		ctx,
		"lock_seats",
		"lock_seats_key",
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        bodyBytes,
			Timestamp:   time.Now(),
		},
	)

	if err != nil {
		return err
	}

	fmt.Println("Published lock seats message in the queue")

	return nil
}

// Unlock seats producer

func (p *Producer) Unlock_Seats(seatsIds []int) error {

	err := p.Conn.ExchangeDeclare(
		"unlock_seats",
		"direct",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	q, err := p.Conn.QueueDeclare(
		"unlock_seats_queue",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	err = p.Conn.QueueBind(
		q.Name,
		"unlock_seats_key",
		"unlock_seats",
		false,
		nil,
	)

	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	bodyBytes, err := json.Marshal(seatsIds)

	if err != nil {
		return err
	}

	err = p.Conn.PublishWithContext(
		ctx,
		"unlock_seats",
		"unlock_seats_key",
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        bodyBytes,
			Timestamp:   time.Now(),
		},
	)

	if err != nil {
		return err
	}

	fmt.Println("Published unlock seats message in the queue")

	return nil
}

// Send email generation

func (p *Producer) Send_Mail_Producer(contactInfo *rabbitmq_producer.Send_Mail_Producer_Request) error {

	err := p.Conn.ExchangeDeclare(
		"send_mail",
		"direct",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	q, err := p.Conn.QueueDeclare(
		"send_mail_queue2",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	err = p.Conn.QueueBind(
		q.Name,
		"send_mail_key",
		"send_mail",
		false,
		nil,
	)

	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	bodyBytes, err := json.Marshal(contactInfo)

	if err != nil {
		return err
	}

	err = p.Conn.PublishWithContext(
		ctx,
		"send_mail",
		"send_mail_key",
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        bodyBytes,
			Timestamp:   time.Now(),
		},
	)

	if err != nil {
		return err
	}

	fmt.Println("Published lock seats message in the queue")

	return nil

}

func (p *Producer) Add_Cast_Producer(cast models.CastAndCrew) error {

	// q, err := p.Conn.QueueDeclare("strapi_create", true, false, false, false, nil)

	// if err != nil {
	// 	return err
	// }

	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	// defer cancel()

	// type Event struct {
	// 	Action string `json:"action"`
	// 	Model  string `json:"model"`
	// 	Data   any    `json:"data"`
	// }

	// castEvent := Event{
	// 	Action: "create",
	// 	Model:  "cast-and-crew",
	// 	Data:   cast,
	// }

	// err = p.Conn.PublishWithContext(
	// 	ctx,
	// 	"",
	// 	q.Name,
	// 	false,
	// 	false,
	// 	amqp091.Publishing{
	// 		ContentType: "application/json",
	// 		Body:        fmt.Appendf(nil, "%v", castEvent),
	// 		Timestamp:   time.Now(),
	// 	},
	// )

	// if err != nil {
	// 	return err
	// }

	// fmt.Println("Published cast creation message in the queue")

	fmt.Println("inside cast and crew producer service")

	return nil
}
