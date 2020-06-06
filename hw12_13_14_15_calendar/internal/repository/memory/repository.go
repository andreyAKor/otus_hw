package memory

import (
	"context"
	"sync"
	"time"

	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/repository"
)

var _ repository.EventsRepo = (*Repo)(nil)

type Repo struct {
	list map[int64]repository.Event // List events by events IDs.
	mx   sync.RWMutex
	id   int64 // Auto-increment id of event record
	once sync.Once
}

// Add new event.
func (r *Repo) Create(ctx context.Context, ev repository.Event) (int64, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.once.Do(r.init)

	d := uint64(ev.Date.Unix() / 60)
	err := r.foreachList(ctx, func(e repository.Event) (err error) {
		if uint64(e.Date.Unix()/60) == d && e.UserID == ev.UserID {
			err = repository.ErrTimeBusy
		}
		return
	})
	if err != nil {
		return 0, err
	}

	r.id++
	ev.ID = r.id
	r.list[ev.ID] = ev

	return ev.ID, nil
}

// Update event by id.
func (r *Repo) Update(ctx context.Context, id int64, ev repository.Event) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.once.Do(r.init)

	if _, ok := r.list[id]; !ok {
		return repository.ErrNotFound
	}

	d := uint64(ev.Date.Unix() / 60)
	err := r.foreachList(ctx, func(e repository.Event) (err error) {
		if uint64(e.Date.Unix()/60) == d && e.UserID == ev.UserID && e.ID != id {
			err = repository.ErrTimeBusy
		}
		return
	})
	if err != nil {
		return err
	}

	r.list[id] = ev

	return nil
}

// Delete event by id.
func (r *Repo) Delete(ctx context.Context, id int64) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	if _, ok := r.list[id]; !ok {
		return repository.ErrNotFound
	}

	delete(r.list, id)

	return nil
}

// Delete last year's events.
func (r *Repo) DeleteOld(ctx context.Context) (err error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	yearAgo := time.Now().AddDate(-1, 0, 0)
	err = r.foreachList(ctx, func(ev repository.Event) error {
		if ev.Date.Unix() <= yearAgo.Unix() {
			delete(r.list, ev.ID)
		}
		return nil
	})
	return
}

// Get list events on date.
func (r *Repo) GetListByDate(ctx context.Context, date time.Time) (events []repository.Event, err error) {
	r.mx.RLock()
	defer r.mx.RUnlock()

	yd := date.YearDay()
	y := date.Year()
	err = r.foreachList(ctx, func(ev repository.Event) error {
		if ev.Date.YearDay() == yd && ev.Date.Year() == y {
			events = append(events, ev)
		}
		return nil
	})
	return
}

// Get list events on week by start date.
func (r *Repo) GetListByWeek(ctx context.Context, start time.Time) (events []repository.Event, err error) {
	r.mx.RLock()
	defer r.mx.RUnlock()

	s := start.Unix()
	e := start.AddDate(0, 0, 7).Unix()
	err = r.foreachList(ctx, func(ev repository.Event) error {
		if ev.Date.Unix() >= s && ev.Date.Unix() <= e {
			events = append(events, ev)
		}
		return nil
	})
	return
}

// Get list events on month by start date.
func (r *Repo) GetListByMonth(ctx context.Context, start time.Time) (events []repository.Event, err error) {
	r.mx.RLock()
	defer r.mx.RUnlock()

	s := start.Unix()
	e := start.AddDate(0, 0, 30).Unix()
	err = r.foreachList(ctx, func(ev repository.Event) error {
		if ev.Date.Unix() >= s && ev.Date.Unix() <= e {
			events = append(events, ev)
		}
		return nil
	})
	return
}

// Init list events as map.
func (r *Repo) init() {
	r.list = make(map[int64]repository.Event)
}

// Iterates list of events and performs the function for each element.
func (r *Repo) foreachList(ctx context.Context, fn func(ev repository.Event) error) error {
	for _, ev := range r.list {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := fn(ev); err != nil {
			return err
		}
	}

	return nil
}
