package pkg

import (
	"github.com/KScaesar/art"
)

type MessageBus interface {
	Send(messages ...*art.Message) error
	RawSend(messages ...*art.Message) error
}

func NewBytesEgress(subject string, bMessage []byte) *art.Message {
	message := art.GetMessage()

	message.Subject = subject
	message.Bytes = bMessage
	return message
}

func NewBodyEgress(subject string, body any) *art.Message {
	message := art.GetMessage()

	message.Subject = subject
	message.Body = body
	return message
}
