package producer

import (
	"context"
	"fmt"
	"time"

	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/rmq"

	"github.com/pkg/errors"
)

const defaultInterval = 5 * time.Second
const contentType = "application/json"

type Producer struct {
	r *rmq.Rmq
}

// Init RabbitMQ producer.
func New(r *rmq.Rmq) (*Producer, error) {
	return &Producer{r}, nil
}

// Running rmq publisher.
func (p *Producer) Run(ctx context.Context) error {
	if err := p.r.Init(ctx); err != nil {
		return errors.Wrap(err, "rmq init fail")
	}

	stop := false
	for !stop {
		select {
		case <-ctx.Done():
			stop = true
		case <-time.After(defaultInterval):
			p.publish()
		}
	}

	return nil
}

func (p *Producer) publish() error {
	fmt.Println("publish after 5sec")

	return p.r.Publish([]byte("lala"), contentType)
}
