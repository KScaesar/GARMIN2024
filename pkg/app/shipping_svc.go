package app

import (
	"context"

	"github.com/KScaesar/art"
)

func NewShippingUseCase() *ShippingUseCase {
	return &ShippingUseCase{}
}

type ShippingUseCase struct {
}

func (uc *ShippingUseCase) CreateShipping(ctx context.Context, event *Order) error {
	shipping := NewShipping(event)
	logger := art.CtxGetLogger(ctx)
	logger.Info("print shipping=%v", art.AnyToString(shipping))
	return nil
}
