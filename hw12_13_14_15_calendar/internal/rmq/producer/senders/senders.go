package senders

import (
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/rmq"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/rmq/producer"
)

const OK = "OK"

type Senders struct {
	producer.Producer
}

// Init RabbitMQ senders-producer.
func New(mq *rmq.Rmq) (*Senders, error) {
	return &Senders{
		Producer: producer.Producer{
			Mq: mq,
		},
	}, nil
}
