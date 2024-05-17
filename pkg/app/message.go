package app

import (
	"github.com/KScaesar/art"
)

var Topic_CreatedOrder = "Created.Order"
var Topic_V1_CreatedOrder = "v1." + Topic_CreatedOrder

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
