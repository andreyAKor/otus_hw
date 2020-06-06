package memory

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/repository"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	ctx := context.Background()

	t.Run("event at this time is busy", func(t *testing.T) {
		r, date := initWithOneEvent(t, ctx)

		_, err := r.Create(ctx, repository.Event{Date: date})
		require.Equal(t, repository.ErrTimeBusy, err)

		_, err = r.Create(ctx, repository.Event{Date: date.Add(time.Microsecond)})
		require.Equal(t, repository.ErrTimeBusy, err)

		_, err = r.Create(ctx, repository.Event{Date: date.Add(time.Minute)})
		require.NoError(t, err)
	})
	t.Run("fail on context cancel", func(t *testing.T) {
		r, date := initWithOneEvent(t, ctx)

		ctx, cancel := context.WithCancel(ctx)
		cancel()

		_, err := r.Create(ctx, repository.Event{Date: date})
		require.Equal(t, context.Canceled, err)
	})
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()

	t.Run("event at this time is busy", func(t *testing.T) {
		r := Repo{}

		date := time.Now()
		ev := repository.Event{Date: date}

		id, err := r.Create(ctx, ev)
		require.NoError(t, err)

		err = r.Update(ctx, id, ev)
		require.NoError(t, err)

		ev.Date = time.Now().AddDate(0, 0, 1)
		_, err = r.Create(ctx, ev)
		require.NoError(t, err)

		err = r.Update(ctx, id, ev)
		require.Equal(t, repository.ErrTimeBusy, err)
	})
	t.Run("fail on context cancel", func(t *testing.T) {
		r := Repo{}

		date := time.Now()
		ev := repository.Event{Date: date}

		id, err := r.Create(ctx, ev)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(ctx)
		cancel()

		err = r.Update(ctx, id, ev)
		require.Equal(t, context.Canceled, err)
	})
	t.Run("event not found", func(t *testing.T) {
		r := Repo{}

		err := r.Update(ctx, 1, repository.Event{})
		require.Equal(t, repository.ErrNotFound, err)
	})
}

func TestDelete(t *testing.T) {
	ctx := context.Background()

	t.Run("normal", func(t *testing.T) {
		r, date := initWithOneEvent(t, ctx)

		id, err := r.Create(ctx, repository.Event{Date: date.Add(time.Minute)})
		require.NoError(t, err)

		err = r.Delete(ctx, id)
		require.NoError(t, err)
	})
	t.Run("event not found", func(t *testing.T) {
		r := Repo{}

		err := r.Delete(ctx, 1)
		require.Equal(t, repository.ErrNotFound, err)
	})
}

func TestDeleteOld(t *testing.T) {
	ctx := context.Background()

	t.Run("normal", func(t *testing.T) {
		r := Repo{}

		date := time.Now().AddDate(-1, 0, 0).Add(-time.Second)

		_, err := r.Create(ctx, repository.Event{Date: date})
		require.NoError(t, err)

		err = r.DeleteOld(ctx)
		require.NoError(t, err)

		_, err = r.Create(ctx, repository.Event{Date: date})
		require.NoError(t, err)
	})
	t.Run("fail on context cancel", func(t *testing.T) {
		r, _ := initWithOneEvent(t, ctx)

		ctx, cancel := context.WithCancel(ctx)
		cancel()

		err := r.DeleteOld(ctx)
		require.Equal(t, context.Canceled, err)
	})
}

func TestGetListByDate(t *testing.T) {
	ctx := context.Background()

	t.Run("normal", func(t *testing.T) {
		r := Repo{}

		date := time.Now()
		ev := repository.Event{ID: 1, Date: date, UserID: 1}

		_, err := r.Create(ctx, ev)
		require.NoError(t, err)

		res, err := r.GetListByDate(ctx, date)
		require.NoError(t, err)
		require.Len(t, res, 1)
		require.Equal(t, ev, res[0])
	})
	t.Run("fail on context cancel", func(t *testing.T) {
		r, date := initWithOneEvent(t, ctx)

		ctx, cancel := context.WithCancel(ctx)
		cancel()

		_, err := r.GetListByDate(ctx, date)
		require.Equal(t, context.Canceled, err)
	})
}

func TestGetListByWeek(t *testing.T) {
	ctx := context.Background()

	t.Run("normal", func(t *testing.T) {
		r := Repo{}

		date := time.Now()
		evList := []repository.Event{
			repository.Event{Date: date.Add(-time.Minute), UserID: 1},
			repository.Event{Date: date, UserID: 1},
			repository.Event{Date: date.AddDate(0, 0, 5), UserID: 1},
			repository.Event{Date: date.AddDate(0, 0, 8), UserID: 1},
		}

		var err error
		for idx, ev := range evList {
			evList[idx].ID, err = r.Create(ctx, ev)
			require.NoError(t, err)
		}

		res, err := r.GetListByWeek(ctx, date)
		require.NoError(t, err)
		require.Len(t, res, 2)
		require.ElementsMatch(t, []repository.Event{evList[1], evList[2]}, res)
	})
	t.Run("fail on context cancel", func(t *testing.T) {
		r, date := initWithOneEvent(t, ctx)

		ctx, cancel := context.WithCancel(ctx)
		cancel()

		_, err := r.GetListByWeek(ctx, date)
		require.Equal(t, context.Canceled, err)
	})
}

func TestForeachList(t *testing.T) {
	ctx := context.Background()

	t.Run("normal", func(t *testing.T) {
		r, _ := initWithOneEvent(t, ctx)

		err := r.foreachList(ctx, func(ev repository.Event) (err error) {
			return
		})
		require.NoError(t, err)
	})
	t.Run("with error", func(t *testing.T) {
		er := errors.New("some error descr")
		r, _ := initWithOneEvent(t, ctx)

		err := r.foreachList(ctx, func(ev repository.Event) (err error) {
			return er
		})
		require.Equal(t, err, er)
	})
	t.Run("fail on context cancel", func(t *testing.T) {
		r, _ := initWithOneEvent(t, ctx)

		ctx, cancel := context.WithCancel(ctx)
		cancel()

		err := r.foreachList(ctx, func(ev repository.Event) (err error) {
			return
		})
		require.Equal(t, context.Canceled, err)
	})
}

func TestMemoryRepoMultithreading(t *testing.T) {
	ctx := context.Background()

	r := Repo{}

	date := time.Now()

	wg := &sync.WaitGroup{}
	wg.Add(7)

	go func() {
		defer wg.Done()
		for i := 0; i < 10_000; i++ {
			_, _ = r.Create(ctx, repository.Event{
				Date: date.Add(time.Hour * time.Duration(i)),
			})
		}
	}()
	go func() {
		defer wg.Done()
		for i := int64(0); i < 10_000; i++ {
			_ = r.Update(ctx, i, repository.Event{
				Date: date.Add(time.Hour*time.Duration(i) + time.Minute),
			})
		}
	}()
	go func() {
		defer wg.Done()
		for i := int64(0); i < 10_000; i++ {
			_ = r.Delete(ctx, i)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 10_000; i++ {
			_ = r.DeleteOld(ctx)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 10_000; i++ {
			_, _ = r.GetListByDate(ctx, date)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 10_000; i++ {
			_, _ = r.GetListByWeek(ctx, date)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 10_000; i++ {
			_, _ = r.GetListByMonth(ctx, date)
		}
	}()

	wg.Wait()
}

func initWithOneEvent(t *testing.T, ctx context.Context) (Repo, time.Time) {
	r := Repo{}

	date := time.Now()

	_, err := r.Create(ctx, repository.Event{Date: date})
	require.NoError(t, err)

	return r, date
}
