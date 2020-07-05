package consumer

import (
	"context"
	"io"
	"sync"

	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/rmq"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

var _ io.Closer = (*Consumer)(nil)

type Consumer struct {
	mq *rmq.Rmq

	consumerTag      string
	qosPrefetchCount int
	threads          int

	done chan struct{}
}

// Init RabbitMQ consumer.
func New(mq *rmq.Rmq, consumerTag string, qosPrefetchCount, threads int) (*Consumer, error) {
	return &Consumer{
		mq,
		consumerTag,
		qosPrefetchCount,
		threads,
		make(chan struct{}),
	}, nil
}

// Running rmq consumer.
func (c *Consumer) Run(ctx context.Context) error {
	if err := c.mq.Init(ctx); err != nil {
		return errors.Wrap(err, "rmq init fail")
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-c.done:
				return
			default:
			}

			if c.mq.IsClosed() {
				continue
			}

			msgsCh, err := c.init()
			if err != nil {
				log.Error().Err(err).Msg("consumer init fail")
				continue
			}

			wg := &sync.WaitGroup{}
			for i := 0; i < c.threads; i++ {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					c.worker(i, msgsCh)
				}(i)
			}
			wg.Wait()
		}
	}()

	return nil
}

func (c *Consumer) Close() error {
	close(c.done)
	return c.mq.Close()
}

func (c *Consumer) worker(workerID int, msgsCh <-chan amqp.Delivery) {
	var conter int

	for msg := range msgsCh {
		conter++
		log.Info().
			Int("workerID", workerID).
			Int("conter", conter).
			Str("msg", string(msg.Body)).
			Msg("message")

		if err := msg.Ack(false); err != nil {
			log.Error().Err(err).Msg("failing acking message of notification")
		}
	}
}

func (c *Consumer) init() (<-chan amqp.Delivery, error) {
	if err := c.mq.Qos(c.qosPrefetchCount); err != nil {
		return nil, errors.Wrap(err, "rmq QOS init fail")
	}

	msgsCh, err := c.mq.Consume(c.consumerTag)
	if err != nil {
		return nil, errors.Wrap(err, "rmq consume init fail")
	}

	return msgsCh, nil
}
