package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	calendarMocks "github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/calendar/mocks"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/repository"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestGetListByMonth(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		t.Run(`"id" parsing fail`, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := calendarMocks.NewMockCalendarer(ctrl)

			srv, err := New(repo, "", 0)
			require.NoError(t, err)

			req := httptest.NewRequest("GET", "/list/month", nil)
			w := httptest.NewRecorder()

			_, err = srv.getListByMonth(w, req)
			require.Error(t, err)

			resp := w.Result()
			require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})
	t.Run("errors", func(t *testing.T) {
		t.Run("can't get events list by month", func(t *testing.T) {
			t.Run("other errors", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				repo := calendarMocks.NewMockCalendarer(ctrl)

				srv, err := New(repo, "", 0)
				require.NoError(t, err)

				now := time.Now()

				req := httptest.NewRequest("GET", "/list/month?start="+now.Format("2006-01-02"), nil)
				w := httptest.NewRecorder()

				tm, err := time.Parse("2006-01-02", now.Format("2006-01-02"))
				require.NoError(t, err)

				repo.EXPECT().
					GetListByMonth(context.Background(), tm).
					Return([]repository.Event{}, errors.New(""))

				_, err = srv.getListByMonth(w, req)
				require.Error(t, err)

				resp := w.Result()
				require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			})
		})
	})
	t.Run("normal", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := calendarMocks.NewMockCalendarer(ctrl)

		srv, err := New(repo, "", 0)
		require.NoError(t, err)

		now := time.Now()

		req := httptest.NewRequest("GET", "/list/month?start="+now.Format("2006-01-02"), nil)
		w := httptest.NewRecorder()

		tm, err := time.Parse("2006-01-02", now.Format("2006-01-02"))
		require.NoError(t, err)

		repo.EXPECT().
			GetListByMonth(context.Background(), tm).
			Return([]repository.Event{}, nil)

		res, err := srv.getListByMonth(w, req)
		require.NoError(t, err)
		require.EqualValues(t, getListByMonth{
			EventList{
				Events: []Event{},
			},
		}, res)

		resp := w.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
