package producer

import (
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/calendar"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/repository"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/rmq"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

const contentType = "application/json"

var _ io.Closer = (*Producer)(nil)

type Producer struct {
	mq       *rmq.Rmq
	calendar calendar.Calendarer

	checkEventsToPublishInterval time.Duration
	checkOldEventsInterval       time.Duration
}

// Init RabbitMQ producer.
func New(
	calendar calendar.Calendarer,
	mq *rmq.Rmq,
	checkEventsToPublishInterval, checkOldEventsInterval string,
) (*Producer, error) {
	checkEventsToPublishIntervalDur, err := time.ParseDuration(checkEventsToPublishInterval)
	if err != nil {
		return nil, errors.Wrapf(err, "check events to publish interval parsing fail (%s)", checkEventsToPublishInterval)
	}

	checkOldEventsIntervalDur, err := time.ParseDuration(checkOldEventsInterval)
	if err != nil {
		return nil, errors.Wrapf(err, "check old events interval parsing fail (%s)", checkOldEventsInterval)
	}

	return &Producer{
		mq,
		calendar,
		checkEventsToPublishIntervalDur,
		checkOldEventsIntervalDur,
	}, nil
}

// Running rmq publisher.
func (p *Producer) Run(ctx context.Context) error {
	if err := p.mq.Init(ctx); err != nil {
		return errors.Wrap(err, "rmq init fail")
	}

	p.worker(ctx, p.checkEventsToPublishInterval, func() error {
		log.Info().
			Str("checkEventsToPublishInterval", p.checkEventsToPublishInterval.String()).
			Msg("checking events to publish")

		if p.mq.IsClosed() {
			return nil
		}

		events, err := p.getEvents(ctx)
		if err != nil {
			return errors.Wrap(err, "failing getting events")
		}

		bodies, err := p.marshaling(events)
		if err != nil {
			return errors.Wrap(err, "failing getting events")
		}

		if err := p.publish(bodies); err != nil {
			log.Error().Err(err).Msg("publishing fail")

			return nil
		}

		return nil
	})
	p.worker(ctx, p.checkOldEventsInterval, func() error {
		log.Info().
			Str("checkOldEventsInterval", p.checkOldEventsInterval.String()).
			Msg("checking old events")

		if err := p.calendar.DeleteOld(ctx); err != nil {
			return errors.Wrap(err, "failing deleting old events")
		}

		return nil
	})

	return nil
}

func (p *Producer) Close() error {
	return p.mq.Close()
}

func (p *Producer) worker(ctx context.Context, d time.Duration, fn func() error) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(d):
				if err := fn(); err != nil {
					log.Fatal().Err(err).Msg("worker error")
					return
				}
			}
		}
	}()
}

// Get list events for publishing.
func (p *Producer) getEvents(ctx context.Context) ([]repository.Event, error) {
	var events []repository.Event

	now := time.Now()

	evList, err := p.calendar.GetListByMonth(ctx, now)
	if err != nil {
		return events, errors.Wrap(err, "failing getting list events on month")
	}

	for _, ev := range evList {
		tmStart := ev.Date
		if ev.DurationStart != nil {
			tmStart = tmStart.Add(-*ev.DurationStart)
		}

		tmEnd := tmStart.Add(ev.Duration)

		if now.After(tmStart) && now.Before(tmEnd) {
			events = append(events, ev)
		}
	}

	return events, nil
}

//nolint:prealloc
// Prepare events to json-string.
func (p *Producer) marshaling(events []repository.Event) ([][]byte, error) {
	var bodies [][]byte

	for _, ev := range events {
		body, err := json.Marshal(rmq.Notification{
			EventID: ev.ID,
			Title:   ev.Title,
			Date:    ev.Date,
			UserID:  ev.UserID,
		})
		if err != nil {
			return bodies, errors.Wrap(err, "json marshaling fail")
		}

		bodies = append(bodies, body)
	}

	return bodies, nil
}

// Publish events list to RabbitMQ.
func (p *Producer) publish(bodies [][]byte) error {
	for _, body := range bodies {
		log.Info().
			Str("body", string(body)).
			Msg("publish")

		if err := p.mq.Publish(amqp.Publishing{
			ContentType: contentType,
			Body:        body,
		}); err != nil {
			return errors.Wrap(err, "rmq publish fail")
		}
	}

	return nil
}
