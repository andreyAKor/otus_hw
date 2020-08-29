package events

import (
	"context"
	"encoding/json"
	"time"

	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/calendar"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/repository"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/rmq"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/rmq/producer"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Events struct {
	producer.ProducerImpl

	calendar calendar.Calendarer

	checkEventsToPublishInterval time.Duration
	checkOldEventsInterval       time.Duration
}

// Init RabbitMQ events-producer.
func New(
	calendar calendar.Calendarer,
	mq *rmq.Rmq,
	checkEventsToPublishInterval, checkOldEventsInterval string,
) (*Events, error) {
	checkEventsToPublishIntervalDur, err := time.ParseDuration(checkEventsToPublishInterval)
	if err != nil {
		return nil, errors.Wrapf(err, "check events to publish interval parsing fail (%s)", checkEventsToPublishInterval)
	}

	checkOldEventsIntervalDur, err := time.ParseDuration(checkOldEventsInterval)
	if err != nil {
		return nil, errors.Wrapf(err, "check old events interval parsing fail (%s)", checkOldEventsInterval)
	}

	return &Events{
		ProducerImpl: producer.ProducerImpl{
			Mq: mq,
		},

		calendar: calendar,

		checkEventsToPublishInterval: checkEventsToPublishIntervalDur,
		checkOldEventsInterval:       checkOldEventsIntervalDur,
	}, nil
}

// Running rmq publisher.
func (e *Events) Run(ctx context.Context) error {
	if err := e.Mq.Init(ctx); err != nil {
		return errors.Wrap(err, "rmq init fail")
	}

	producer.Worker(ctx, e.checkEventsToPublishInterval, func(ctx context.Context) error {
		log.Info().
			Str("checkEventsToPublishInterval", e.checkEventsToPublishInterval.String()).
			Msg("checking events to publish")

		if e.Mq.IsClosed() {
			return nil
		}

		events, err := e.getEvents(ctx)
		if err != nil {
			return errors.Wrap(err, "failing getting events")
		}

		bodies, err := e.marshaling(events)
		if err != nil {
			return errors.Wrap(err, "failing getting events")
		}

		if err := e.publish(bodies); err != nil {
			log.Error().
				Err(err).
				Msg("publishing fail")

			return nil
		}

		return nil
	})
	producer.Worker(ctx, e.checkOldEventsInterval, func(ctx context.Context) error {
		log.Info().
			Str("checkOldEventsInterval", e.checkOldEventsInterval.String()).
			Msg("checking old events")

		if err := e.calendar.DeleteOld(ctx); err != nil {
			return errors.Wrap(err, "failing deleting old events")
		}

		return nil
	})

	return nil
}

// Get list events for publishing.
func (e *Events) getEvents(ctx context.Context) ([]repository.Event, error) {
	var events []repository.Event

	now := time.Now()

	evList, err := e.calendar.GetListByMonth(ctx, now)
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
func (e *Events) marshaling(events []repository.Event) ([][]byte, error) {
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
func (e *Events) publish(bodies [][]byte) error {
	for _, body := range bodies {
		if err := e.Publish(body); err != nil {
			return errors.Wrap(err, "producer publish fail")
		}
	}

	return nil
}
