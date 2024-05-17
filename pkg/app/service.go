package app

type Service struct {
	OrderService    *OrderUseCase
	ShippingService *ShippingUseCase
}
