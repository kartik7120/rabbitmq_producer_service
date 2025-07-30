package producers

import (
	"context"
	"errors"
	"fmt"
	"time"

	rabbitmq_producer "github.com/kartik7120/booking_rabbitmq_producer_service/cmd/grpcServer"
	"github.com/kartik7120/booking_rabbitmq_producer_service/cmd/models"
)

type Rabbitmq_Producer_Service struct {
	rabbitmq_producer.UnimplementedRabbitmqProducerServiceServer
	Producer Producer
}

func NewRabbitmq_Producer_Service() *Rabbitmq_Producer_Service {
	return &Rabbitmq_Producer_Service{}
}

func (r *Rabbitmq_Producer_Service) Payment_Service_Webhook_Producer(ctx context.Context, in *rabbitmq_producer.Payment_Service_Producer_Request) (*rabbitmq_producer.Payment_Service_Producer_Response, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	done := make(chan error, 1)

	go func() {
		var requestPayload models.Payment

		// Billing
		requestPayload.Billing.City = in.PaymentPayload.Billing.City
		requestPayload.Billing.Country = in.PaymentPayload.Billing.Country
		requestPayload.Billing.State = in.PaymentPayload.Billing.State
		requestPayload.Billing.Street = in.PaymentPayload.Billing.Street
		requestPayload.Billing.Zipcode = in.PaymentPayload.Billing.Zipcode

		// General Info
		requestPayload.BrandID = in.PaymentPayload.BrandId
		requestPayload.BusinessID = in.PaymentPayload.BusinessId
		requestPayload.CardIssuingCountry = &in.PaymentPayload.CardIssuingCountry
		requestPayload.CardLastFour = in.PaymentPayload.CardLastFour
		requestPayload.CardNetwork = in.PaymentPayload.CardNetwork
		requestPayload.CardType = in.PaymentPayload.CardType
		requestPayload.CreatedAt = in.PaymentPayload.CreatedAt.AsTime()
		requestPayload.UpdatedAt = in.PaymentPayload.UpdatedAt.AsTime()
		requestPayload.Currency = in.PaymentPayload.Currency
		requestPayload.PaymentID = in.PaymentPayload.PaymentId
		requestPayload.PaymentLink = in.PaymentPayload.PaymentLink
		requestPayload.PaymentMethod = in.PaymentPayload.PaymentMethod
		requestPayload.PaymentMethodType = in.PaymentPayload.PaymentMethodType
		requestPayload.Status = &in.PaymentPayload.Status
		requestPayload.SubscriptionID = in.PaymentPayload.SubscriptionId
		requestPayload.Tax = int(in.PaymentPayload.Tax)
		requestPayload.TotalAmount = int(in.PaymentPayload.TotalAmount)
		requestPayload.SettlementAmount = int(in.PaymentPayload.SettlementAmount)
		requestPayload.SettlementCurrency = in.PaymentPayload.SettlementCurrency
		requestPayload.SettlementTax = int(in.PaymentPayload.SettlementTax)
		requestPayload.DigitalProductsDelivered = in.PaymentPayload.DigitalProductsDelivered
		requestPayload.DiscountID = in.PaymentPayload.DiscountId
		requestPayload.ErrorCode = in.PaymentPayload.ErrorCode
		requestPayload.ErrorMessage = in.PaymentPayload.ErrorMessage

		// Customer
		requestPayload.Customer.CustomerID = in.PaymentPayload.Customer.CustomerId
		requestPayload.Customer.Email = in.PaymentPayload.Customer.Email
		requestPayload.Customer.Name = in.PaymentPayload.Customer.Name

		// Product Cart
		for _, p := range in.PaymentPayload.ProductCart {
			requestPayload.ProductCart = append(requestPayload.ProductCart, struct {
				ProductID string `json:"product_id"`
				Quantity  int    `json:"quantity"`
			}{
				ProductID: p.ProductId,
				Quantity:  int(p.Quantity),
			})
		}

		// Disputes
		for _, d := range in.PaymentPayload.Disputes {
			requestPayload.Disputes = append(requestPayload.Disputes, struct {
				Amount        string    `json:"amount"`
				BusinessID    string    `json:"business_id"`
				CreatedAt     time.Time `json:"created_at"`
				Currency      string    `json:"currency"`
				DisputeID     string    `json:"dispute_id"`
				DisputeStage  string    `json:"dispute_stage"`
				DisputeStatus string    `json:"dispute_status"`
				PaymentID     string    `json:"payment_id"`
				Remarks       string    `json:"remarks"`
			}{
				Amount:        d.Amount,
				BusinessID:    d.BusinessId,
				CreatedAt:     d.CreatedAt.AsTime(),
				Currency:      d.Currency,
				DisputeID:     d.DisputeId,
				DisputeStage:  d.DisputeStage,
				DisputeStatus: d.DisputeStatus,
				PaymentID:     d.PaymentId,
				Remarks:       d.Remarks,
			})
		}

		// Refunds
		for _, rf := range in.PaymentPayload.Refunds {
			requestPayload.Refunds = append(requestPayload.Refunds, struct {
				Amount     int       `json:"amount"`
				BusinessID string    `json:"business_id"`
				CreatedAt  time.Time `json:"created_at"`
				Currency   *string   `json:"currency"`
				IsPartial  bool      `json:"is_partial"`
				PaymentID  string    `json:"payment_id"`
				Reason     string    `json:"reason"`
				RefundID   string    `json:"refund_id"`
				Status     string    `json:"status"`
			}{
				Amount:     int(rf.Amount),
				BusinessID: rf.BusinessId,
				CreatedAt:  rf.CreatedAt.AsTime(),
				Currency:   &rf.Currency,
				IsPartial:  rf.IsPartial,
				PaymentID:  rf.PaymentId,
				Reason:     rf.Reason,
				RefundID:   rf.RefundId,
				Status:     rf.Status,
			})
		}

		// Metadata (if provided as JSON string map)
		// if in.PaymentPayload.Metadata != nil {
		// 	requestPayload.Metadata = make(map[string]interface{})
		// 	for key, val := range in.PaymentPayload.Metadata.Fields {
		// 		// Simplified conversion
		// 		requestPayload.Metadata[key] = val.GetStringValue()
		// 	}
		// }

		// Send to RabbitMQ
		err := r.Producer.Payment_Service_Producer(requestPayload)

		done <- err
	}()

	select {
	case <-ctx.Done():
		fmt.Println("Timeout as request took more than 10 seconds")
		return &rabbitmq_producer.Payment_Service_Producer_Response{
			Error: ctx.Err().Error(),
		}, ctx.Err()
	case err := <-done:
		if err != nil {
			return nil, err
		}
	}

	return &rabbitmq_producer.Payment_Service_Producer_Response{
		Error: "",
	}, nil
}

func (r *Rabbitmq_Producer_Service) Lock_Seats(ctx context.Context, in *rabbitmq_producer.Lock_Seats_Request) (*rabbitmq_producer.Lock_Seats_Response, error) {

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	done := make(chan error, 1)

	go func() {
		var seatIds []int

		for _, v := range in.SeatIds {
			seatIds = append(seatIds, int(v))
		}

		err := r.Producer.Lock_Seats(seatIds)

		done <- err
	}()

	select {
	case <-ctx.Done():
		fmt.Println("Timeout 10 seconds reached")
	case err := <-done:

		if err != nil {
			return nil, err
		}
	}

	return &rabbitmq_producer.Lock_Seats_Response{
		Error: "",
	}, nil
}

func (r *Rabbitmq_Producer_Service) Unlock_Seats(ctx context.Context, in *rabbitmq_producer.Unlock_Seats_Request) (*rabbitmq_producer.Unlock_Seats_Response, error) {

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	done := make(chan error, 1)

	go func() {
		var seatIds []int

		for _, v := range in.SeatIds {
			seatIds = append(seatIds, int(v))
		}

		err := r.Producer.Unlock_Seats(seatIds)

		done <- err
	}()

	select {
	case <-ctx.Done():
		fmt.Println("Timeout 10 second")
		return nil, errors.New("timeout 10 second exceeded for sending Unlock seats into the queue")
	case err := <-done:
		if err != nil {
			return nil, err
		}
	}

	return &rabbitmq_producer.Unlock_Seats_Response{
		Error: "",
	}, nil
}

func (r *Rabbitmq_Producer_Service) Send_Mail_Producer(ctx context.Context, in *rabbitmq_producer.Send_Mail_Producer_Request) (*rabbitmq_producer.Send_Mail_Producer_Response, error) {

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	done := make(chan error, 1)

	go func() {
		err := r.Producer.Send_Mail_Producer(struct {
			Email        string
			Phone_number string
		}{
			Email:        in.Email,
			Phone_number: in.PhoneNumber,
		})
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			return nil, err
		}
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	return &rabbitmq_producer.Send_Mail_Producer_Response{
		Error: "",
	}, nil
}
