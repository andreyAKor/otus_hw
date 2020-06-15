package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/repository"
	repositoryMocks "github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/repository/mocks"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestDelete(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		t.Run(`"id" parsing fail`, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := repositoryMocks.NewMockEventsRepo(ctrl)

			srv, err := New(repo, "", 0)
			require.NoError(t, err)

			req := httptest.NewRequest("DELETE", "/delete", nil)
			w := httptest.NewRecorder()

			_, err = srv.delete(w, req)
			require.Error(t, err)

			resp := w.Result()
			require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})
	t.Run("errors", func(t *testing.T) {
		t.Run("can't delete event", func(t *testing.T) {
			t.Run("event not found", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				repo := repositoryMocks.NewMockEventsRepo(ctrl)

				srv, err := New(repo, "", 0)
				require.NoError(t, err)

				req := httptest.NewRequest("DELETE", "/delete?id=0", nil)
				w := httptest.NewRecorder()

				repo.EXPECT().Delete(context.Background(), int64(0)).Return(repository.ErrNotFound)

				_, err = srv.delete(w, req)
				require.Equal(t, repository.ErrNotFound, errors.Cause(err))

				resp := w.Result()
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			})
			t.Run("other errors", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				repo := repositoryMocks.NewMockEventsRepo(ctrl)

				srv, err := New(repo, "", 0)
				require.NoError(t, err)

				req := httptest.NewRequest("DELETE", "/delete?id=0", nil)
				w := httptest.NewRecorder()

				repo.EXPECT().Delete(context.Background(), int64(0)).Return(errors.New(""))

				_, err = srv.delete(w, req)
				require.Error(t, err)

				resp := w.Result()
				require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			})
		})
	})
	t.Run("normal", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositoryMocks.NewMockEventsRepo(ctrl)

		srv, err := New(repo, "", 0)
		require.NoError(t, err)

		req := httptest.NewRequest("DELETE", "/delete?id=0", nil)
		w := httptest.NewRecorder()

		repo.EXPECT().Delete(context.Background(), int64(0)).Return(nil)

		res, err := srv.delete(w, req)
		require.NoError(t, err)
		require.Nil(t, res)

		resp := w.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
