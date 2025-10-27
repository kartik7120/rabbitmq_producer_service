package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	"time"

	rabbitmq_producer "github.com/kartik7120/booking_rabbitmq_producer_service/cmd/grpcServer"
	"github.com/kartik7120/booking_rabbitmq_producer_service/cmd/producers"
	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func connectRabbitMQ(url string, retries int, delay time.Duration) (*amqp091.Connection, error) {
	var conn *amqp091.Connection
	var err error

	for i := 0; i < retries; i++ {
		conn, err = amqp091.Dial(url)
		if err == nil {
			return conn, nil
		}
		fmt.Printf("Retrying RabbitMQ connection (%d/%d)...\n", i+1, retries)
		time.Sleep(delay)
	}
	return nil, err
}

func main() {
	fmt.Println("RabbitMQ Producer Service is running...")

	client, err := connectRabbitMQ("amqp://guest:guest@rabbitmq_booking_app:5672/", 10, 3*time.Second)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	if err != nil {
		fmt.Printf("Failed to connect to RabbitMQ: %s\n", err)
		os.Exit(1)
		return
	}

	defer client.Close()

	channel, err := client.Channel()

	if err != nil {
		fmt.Printf("Failed to open a channel: %s\n", err)
		os.Exit(1)
		return
	}

	defer channel.Close()

	var opts []grpc.ServerOption

	lis, err := net.Listen("tcp", ":1105")

	server := grpc.NewServer(opts...)

	rabbitmq_producer.RegisterRabbitmqProducerServiceServer(
		server, &producers.Rabbitmq_Producer_Service{
			Producer: producers.Producer{
				Conn: channel,
			},
		},
	)

	if os.Getenv("ENV") != "production" {
		reflection.Register(server)
	}

	go func() {
		if err := server.Serve(lis); err != nil {
			fmt.Println("error starting the rabbitmq producer service")
			panic(err)
		}
	}()

	if err != nil {
		fmt.Printf("Failed to listen on port 1105: %s\n", err)
		os.Exit(1)
		return
	}

	<-ch

	fmt.Println("Shutting down RabbitMQ Producer Service gracefully...")

	server.GracefulStop()

}
