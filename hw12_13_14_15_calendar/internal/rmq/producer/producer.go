package producer

import (
	"context"

	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/rmq"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

const contentType = "application/json"

var _ Producer = (*ProducerImpl)(nil)

//nolint:golint
type ProducerImpl struct {
	Mq *rmq.Rmq
}

// Running rmq publisher.
func (p *ProducerImpl) Run(ctx context.Context) error {
	if err := p.Mq.Init(ctx); err != nil {
		return errors.Wrap(err, "rmq init fail")
	}

	return nil
}

// Publish content to RabbitMQ.
func (p *ProducerImpl) Publish(content []byte) error {
	log.Info().
		Str("content", string(content)).
		Msg("publish")

	if err := p.Mq.Publish(amqp.Publishing{
		ContentType: contentType,
		Body:        content,
	}); err != nil {
		return errors.Wrap(err, "rmq publish fail")
	}

	return nil
}

func (p *ProducerImpl) Close() error {
	return p.Mq.Close()
}
