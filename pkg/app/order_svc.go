package app

import (
	"context"

	"github.com/KScaesar/GARMIN2024/pkg"
)

func NewOrderUseCase(bus pkg.MessageBus) *OrderUseCase {
	return &OrderUseCase{bus: bus}
}

type OrderUseCase struct {
	bus pkg.MessageBus
}

func (uc *OrderUseCase) CreateOrder(ctx context.Context, req *CreateOrderParam) error {
	_, event := NewOrder(req)

	err := uc.bus.Send(event)
	if err != nil {
		return err
	}
	return nil
}
