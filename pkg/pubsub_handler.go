package pkg

import (
	"encoding/json"

	"github.com/KScaesar/art"
)

func UseEncodeJson() art.Middleware {
	return func(next art.HandleFunc) art.HandleFunc {
		return func(message *art.Message, dependency any) error {
			if message.Bytes != nil {
				return nil
			}

			message.Mutex.Lock()

			if message.Bytes != nil {
				message.Mutex.Unlock()
				return nil
			}

			bytes, err := json.Marshal(message.Body)
			if err != nil {
				message.Mutex.Unlock()
				return err
			}

			message.Bytes = bytes
			message.Mutex.Unlock()
			return next(message, dependency)
		}
	}
}
