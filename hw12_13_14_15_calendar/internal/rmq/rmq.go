package rmq

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

var (
	ErrStopReconn = errors.New("stop reconnecting")
)

type Rmq struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	done         chan error
	uri          string
	exchangeName string
	exchangeType string
	queueName    string
	bindingKey   string
}

// RabbitMQ connector.
func New(uri, exchangeName, exchangeType, queueName, bindingKey string) (*Rmq, error) {
	return &Rmq{
		done:         make(chan error),
		uri:          uri,
		exchangeName: exchangeName,
		exchangeType: exchangeType,
		queueName:    queueName,
		bindingKey:   bindingKey,
	}, nil
}

// Init RabbitMQ.
func (r *Rmq) Init(ctx context.Context) error {
	if err := r.connect(); err != nil {
		return errors.Wrap(err, "rmq connection fail")
	}

	if err := r.prepareQueue(); err != nil {
		return errors.Wrap(err, "announce queue fail")
	}

	return nil
}

// Reconnecting algo
func (r *Rmq) reConnect() error {
	be := backoff.NewExponentialBackOff()
	be.MaxElapsedTime = time.Minute
	be.InitialInterval = 1 * time.Second
	be.Multiplier = 2
	be.MaxInterval = 15 * time.Second

	b := backoff.WithContext(be, context.Background())
	for {
		d := b.NextBackOff()
		if d == backoff.Stop {
			return ErrStopReconn
		}

		select {
		case <-time.After(d):
			if err := r.connect(); err != nil {
				log.Error().Err(err).Msg("could not connect in reconnect call")
				continue
			}
			if err := r.prepareQueue(); err != nil {
				log.Error().Err(err).Msg("couldn't connect")
				continue
			}

			return nil
		}
	}
}

// Connect to RabbitMQ.
func (r *Rmq) connect() error {
	var err error

	r.conn, err = amqp.Dial(r.uri)
	if err != nil {
		return errors.Wrap(err, "dial fail")
	}

	r.channel, err = r.conn.Channel()
	if err != nil {
		return errors.Wrap(err, "channel fail")
	}

	go func() {
		log.Printf("closing: %s", <-r.conn.NotifyClose(make(chan *amqp.Error)))
		// Понимаем, что канал сообщений закрыт, надо пересоздать соединение.
		r.done <- errors.New("Channel Closed")
	}()

	if err := r.channel.ExchangeDeclare(
		r.exchangeName,
		r.exchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return errors.Wrap(err, "exchange declare fail")
	}

	return nil
}

// Declare queue.
func (r *Rmq) prepareQueue() error {
	_, err := r.channel.QueueDeclare(
		r.queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "queue declare fail")
	}

	/*
		// Число сообщений, которые можно подтвердить за раз.
		err = p.channel.Qos(50, 0, false)
		if err != nil {
			return nil, fmt.Errorf("Error setting qos: %s", err)
		}
	*/

	// Создаём биндинг (правило маршрутизации).
	if err = r.channel.QueueBind(
		r.queueName,
		r.bindingKey,
		r.exchangeName,
		false,
		nil,
	); err != nil {
		return errors.Wrap(err, "queue bind fail")
	}

	/*
		msgs, err := p.channel.Consume(
			r.queueName,
			c.consumerTag,
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("Queue Consume: %s", err)
		}

		return msgs, nil
		/**/

	return nil
}

func (r *Rmq) Publish(body []byte, ct string) error {
	return r.channel.Publish(
		r.exchangeName, // exchange
		r.queueName,    // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: ct,
			Body:        body,
		})
}
