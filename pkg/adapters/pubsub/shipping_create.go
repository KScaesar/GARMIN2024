package pubsub

import (
	"encoding/json"

	"github.com/KScaesar/art"

	"github.com/KScaesar/GARMIN2024/pkg/app"
)

func createShipping(svc *app.Service) art.HandleFunc {
	return func(ingress *art.Message, dep any) error {
		var event app.Order
		err := json.Unmarshal(ingress.Bytes, &event)
		if err != nil {
			return err
		}
		return svc.ShippingService.CreateShipping(ingress.Ctx, &event)
	}
}
