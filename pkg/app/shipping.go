package app

import (
	"github.com/KScaesar/art"
)

// domain entity

func NewShipping(event *Order) *Shipping {
	shipping := &Shipping{
		ShippingId:   art.GenerateRandomCode(8),
		OrderId:      event.OrderId,
		CustomerName: event.CustomerName,
	}
	return shipping
}

type Shipping struct {
	ShippingId   string `json:"shipping_id"`
	OrderId      string `json:"order_id"`
	CustomerName string `json:"customer_name"`
}
