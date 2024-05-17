package pkg

import (
	"encoding/json"

	"github.com/KScaesar/art"
)

func UseEncodeJson() art.Middleware {
	return func(next art.HandleFunc) art.HandleFunc {
		return func(message *art.Message, dependency any) error {
			bytes, err := json.Marshal(message.Body)
			if err != nil {
				return err
			}
			message.Bytes = bytes
			return next(message, dependency)
		}
	}
}
