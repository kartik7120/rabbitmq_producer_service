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

	// 1. Ensure exchange exists (idempotent)
	err := p.Conn.ExchangeDeclare(
		"send_mail",
		"direct",
		true, // durable
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// 2. Marshal payload
	bodyBytes, err := json.Marshal(contactInfo)
	if err != nil {
		return err
	}

	// Optional: generate a unique ID for tracking
	// 3. Publish
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = p.Conn.PublishWithContext(
		ctx,
		"send_mail",     // exchange
		"send_mail_key", // routing key
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

	fmt.Println("Published send_mail message")

	return nil
}

func (p *Producer) Add_Cast_Producer(cast ExtendedCastAndCrew) error {

	// Declare queue
	q, err := p.Conn.QueueDeclare(
		"strapi_create",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("queue declare failed: %w", err)
	}

	// Declare exchange
	err = p.Conn.ExchangeDeclare(
		"strapi_create_exchange",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("exchange declare failed: %w", err)
	}

	// Bind queue → exchange
	err = p.Conn.QueueBind(
		q.Name,
		"cast_creation", // routing key
		"strapi_create_exchange",
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("queue bind failed: %w", err)
	}

	// Prepare message
	castEvent := struct {
		Action string `json:"action"`
		Model  string `json:"model"`
		Data   any    `json:"data"`
	}{
		Action: "create",
		Model:  "cast-and-crew",
		Data:   cast,
	}

	payload, err := json.Marshal(castEvent)
	if err != nil {
		return err
	}

	// Publish with unique ID
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = p.Conn.PublishWithContext(
		ctx,
		"strapi_create_exchange",
		"cast_creation",
		false,
		false,
		amqp091.Publishing{
			ContentType:   "application/json",
			Body:          payload,
			Timestamp:     time.Now(),
			MessageId:     cast.StarpiCastUid,
			CorrelationId: cast.StarpiCastUid,
		},
	)
	if err != nil {
		return fmt.Errorf("publish failed: %w", err)
	}

	fmt.Println("Published cast creation message in the queue")
	return nil
}

func (p *Producer) Delete_Cast_Producer(cast ExtendedCastAndCrew) error {
	q, err := p.Conn.QueueDeclare(
		"strapi_create",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("queue declare failed: %w", err)
	}

	err = p.Conn.ExchangeDeclare(
		"strapi_create_exchange",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("exchange declare failed: %w", err)
	}

	err = p.Conn.QueueBind(
		q.Name,
		"cast_deletion",
		"strapi_create_exchange",
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("queue bind failed: %w", err)
	}

	payload := struct {
		Action string `json:"action"`
		Model  string `json:"model"`
		Data   any    `json:"data"`
	}{
		Action: "delete",
		Model:  "cast-and-crew",
		Data:   cast,
	}

	body, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = p.Conn.PublishWithContext(
		ctx,
		"strapi_create_exchange",
		"cast_deletion",
		false,
		false,
		amqp091.Publishing{
			ContentType:   "application/json",
			Body:          body,
			Timestamp:     time.Now(),
			MessageId:     cast.StarpiCastUid,
			CorrelationId: cast.StarpiCastUid,
		},
	)

	if err != nil {
		return fmt.Errorf("publish failed: %w", err)
	}

	fmt.Println("Published cast deletion message in the queue")
	return nil
}

func (p *Producer) Movie_Time_Slot_Producer(payload MovieTimeSlotPayload) error {

	q, err := p.Conn.QueueDeclare(
		"strapi_create",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("queue declare failed: %w", err)
	}

	err = p.Conn.ExchangeDeclare(
		"strapi_create_exchange",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("exchange declare failed: %w", err)
	}

	err = p.Conn.QueueBind(
		q.Name,
		"movie_time_slot_creation",
		"strapi_create_exchange",
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("queue bind failed: %w", err)
	}

	payloadEvent := struct {
		Action string `json:"action"`
		Model  string `json:"model"`
		Data   any    `json:"data"`
	}{
		Action: "create",
		Model:  "movie-time-slot",
		Data:   payload,
	}

	body, err := json.Marshal(payloadEvent)

	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = p.Conn.PublishWithContext(
		ctx,
		"strapi_create_exchange",
		"movie_time_slot_creation",
		false,
		false,
		amqp091.Publishing{
			ContentType:   "application/json",
			Body:          body,
			Timestamp:     time.Now(),
			MessageId:     payload.StarpiMovieUid,
			CorrelationId: payload.StarpiMovieUid,
		},
	)

	if err != nil {
		return fmt.Errorf("publish failed: %w", err)
	}

	fmt.Println("Published movie time slot creation message in the queue")
	return nil
}
