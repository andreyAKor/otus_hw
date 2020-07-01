package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	calendarMocks "github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/calendar/mocks"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/repository"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	postForm := url.Values{
		"title":    {"title-123"},
		"date":     {"2018-09-22T19:42:31+03:00"},
		"duration": {"0h"},
		"user_id":  {"32"},
	}

	t.Run("empty", func(t *testing.T) {
		t.Run("can't preparing event", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := calendarMocks.NewMockCalendarer(ctrl)

			srv, err := New(repo, "", 0)
			require.NoError(t, err)

			req := httptest.NewRequest("POST", "/create", nil)
			w := httptest.NewRecorder()

			_, err = srv.create(w, req)
			require.Error(t, err)

			resp := w.Result()
			require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})
	t.Run("errors", func(t *testing.T) {
		t.Run("can't create new event", func(t *testing.T) {
			t.Run("event at this time is busy", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				repo := calendarMocks.NewMockCalendarer(ctrl)

				srv, err := New(repo, "", 0)
				require.NoError(t, err)

				req := httptest.NewRequest("POST", "/create", nil)
				req.PostForm = postForm
				w := httptest.NewRecorder()

				ev, err := srv.prepareEventRequest(req)
				require.NoError(t, err)

				repo.EXPECT().Create(context.Background(), *ev).Return(int64(0), repository.ErrTimeBusy)

				_, err = srv.create(w, req)
				require.Equal(t, repository.ErrTimeBusy, errors.Cause(err))

				resp := w.Result()
				require.Equal(t, http.StatusConflict, resp.StatusCode)
			})
			t.Run("other errors", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				repo := calendarMocks.NewMockCalendarer(ctrl)

				srv, err := New(repo, "", 0)
				require.NoError(t, err)

				req := httptest.NewRequest("POST", "/create", nil)
				req.PostForm = postForm
				w := httptest.NewRecorder()

				ev, err := srv.prepareEventRequest(req)
				require.NoError(t, err)

				repo.EXPECT().Create(context.Background(), *ev).Return(int64(0), errors.New(""))

				_, err = srv.create(w, req)
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

		req := httptest.NewRequest("POST", "/create", nil)
		req.PostForm = postForm
		w := httptest.NewRecorder()

		ev, err := srv.prepareEventRequest(req)
		require.NoError(t, err)

		id := int64(1)

		repo.EXPECT().Create(context.Background(), *ev).Return(id, nil)

		res, err := srv.create(w, req)
		require.NoError(t, err)
		require.EqualValues(t, create{id}, res)

		resp := w.Result()
		require.Equal(t, http.StatusCreated, resp.StatusCode)
	})
}
