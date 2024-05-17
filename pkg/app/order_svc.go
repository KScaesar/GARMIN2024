package app

import (
	"context"
)

func NewOrderUseCase(bus MessageBus) *OrderUseCase {
	return &OrderUseCase{bus: bus}
}

type OrderUseCase struct {
	bus MessageBus
}

func (uc *OrderUseCase) CreateOrder(ctx context.Context, req *CreateOrderParam) error {
	_, event := NewOrder(req)

	err := uc.bus.Send(event)
	if err != nil {
		return err
	}
	return nil
}
