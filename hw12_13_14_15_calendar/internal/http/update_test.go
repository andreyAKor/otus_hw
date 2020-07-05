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

func TestUpdate(t *testing.T) {
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

			req := httptest.NewRequest("PUT", "/update", nil)
			w := httptest.NewRecorder()

			_, err = srv.update(w, req)
			require.Error(t, err)

			resp := w.Result()
			require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
		t.Run(`"id" parsing fail`, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := calendarMocks.NewMockCalendarer(ctrl)

			srv, err := New(repo, "", 0)
			require.NoError(t, err)

			req := httptest.NewRequest("PUT", "/update", nil)
			w := httptest.NewRecorder()

			_, err = srv.update(w, req)
			require.Error(t, err)

			resp := w.Result()
			require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})
	t.Run("errors", func(t *testing.T) {
		t.Run("can't update event", func(t *testing.T) {
			t.Run("event at this time is busy", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				repo := calendarMocks.NewMockCalendarer(ctrl)

				srv, err := New(repo, "", 0)
				require.NoError(t, err)

				req := httptest.NewRequest("PUT", "/update?id=0", nil)
				req.PostForm = postForm
				w := httptest.NewRecorder()

				ev, err := srv.prepareEventRequest(req)
				require.NoError(t, err)

				repo.EXPECT().Update(context.Background(), int64(0), *ev).Return(repository.ErrTimeBusy)

				_, err = srv.update(w, req)
				require.Equal(t, repository.ErrTimeBusy, errors.Cause(err))

				resp := w.Result()
				require.Equal(t, http.StatusConflict, resp.StatusCode)
			})
			t.Run("event not found", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				repo := calendarMocks.NewMockCalendarer(ctrl)

				srv, err := New(repo, "", 0)
				require.NoError(t, err)

				req := httptest.NewRequest("PUT", "/update?id=0", nil)
				req.PostForm = postForm
				w := httptest.NewRecorder()

				ev, err := srv.prepareEventRequest(req)
				require.NoError(t, err)

				repo.EXPECT().Update(context.Background(), int64(0), *ev).Return(repository.ErrNotFound)

				_, err = srv.update(w, req)
				require.Equal(t, repository.ErrNotFound, errors.Cause(err))

				resp := w.Result()
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			})
			t.Run("other errors", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				repo := calendarMocks.NewMockCalendarer(ctrl)

				srv, err := New(repo, "", 0)
				require.NoError(t, err)

				req := httptest.NewRequest("PUT", "/update?id=0", nil)
				req.PostForm = postForm
				w := httptest.NewRecorder()

				ev, err := srv.prepareEventRequest(req)
				require.NoError(t, err)

				repo.EXPECT().Update(context.Background(), int64(0), *ev).Return(errors.New(""))

				_, err = srv.update(w, req)
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

		req := httptest.NewRequest("PUT", "/update?id=0", nil)
		req.PostForm = postForm
		w := httptest.NewRecorder()

		ev, err := srv.prepareEventRequest(req)
		require.NoError(t, err)

		repo.EXPECT().Update(context.Background(), int64(0), *ev).Return(nil)

		res, err := srv.update(w, req)
		require.NoError(t, err)
		require.Nil(t, res)

		resp := w.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
