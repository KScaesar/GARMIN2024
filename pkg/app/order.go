package app

import (
	"github.com/KScaesar/art"

	"github.com/KScaesar/GARMIN2024/pkg"
)

// domain entity

type CreateOrderParam struct {
	CustomerName string `json:"customer_name"`
	TotalPrice   int    `json:"total_price"`
}

func NewOrder(param *CreateOrderParam) (order *Order, event *art.Message) {
	order = &Order{
		OrderId:      art.GenerateRandomCode(6),
		CustomerName: param.CustomerName,
		TotalPrice:   param.TotalPrice,
	}

	event = pkg.NewBodyEgress(Subject_V1_CreatedOrder, order)
	Metadata.SetOrderId(event, order.OrderId)

	return order, event
}

type Order struct {
	OrderId      string `json:"order_id"`
	CustomerName string `json:"customer_name"`
	TotalPrice   int    `json:"total_price"`
}
