package grpc

import (
	"context"
	"testing"

	calendarMocks "github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/calendar/mocks"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/repository"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/schema"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestDelete(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		t.Run("event not found", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := calendarMocks.NewMockCalendarer(ctrl)

			srv, err := New(repo, "", 0)
			require.NoError(t, err)

			repo.EXPECT().Delete(context.Background(), int64(0)).Return(repository.ErrNotFound)

			resp, err := srv.Delete(context.Background(), &schema.DeleteRpcRequest{})
			require.Equal(t, repository.ErrNotFound, errors.Cause(err))
			require.Nil(t, resp)
		})
	})
	t.Run("errors", func(t *testing.T) {
		t.Run("can't create new event", func(t *testing.T) {
			t.Run("event not found", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				repo := calendarMocks.NewMockCalendarer(ctrl)

				srv, err := New(repo, "", 0)
				require.NoError(t, err)

				repo.EXPECT().Delete(context.Background(), int64(0)).Return(repository.ErrNotFound)

				resp, err := srv.Delete(context.Background(), &schema.DeleteRpcRequest{Id: 0})
				require.Equal(t, repository.ErrNotFound, errors.Cause(err))
				require.Nil(t, resp)
			})
			t.Run("other errors", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				repo := calendarMocks.NewMockCalendarer(ctrl)

				srv, err := New(repo, "", 0)
				require.NoError(t, err)

				repo.EXPECT().Delete(context.Background(), int64(0)).Return(errors.New(""))

				resp, err := srv.Delete(context.Background(), &schema.DeleteRpcRequest{Id: 0})
				require.Error(t, err)
				require.Nil(t, resp)
			})
		})
	})
	t.Run("normal", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := calendarMocks.NewMockCalendarer(ctrl)

		srv, err := New(repo, "", 0)
		require.NoError(t, err)

		repo.EXPECT().Delete(context.Background(), int64(0)).Return(nil)

		resp, err := srv.Delete(context.Background(), &schema.DeleteRpcRequest{Id: 0})
		require.NoError(t, err)
		require.EqualValues(t, empty.Empty{}, *resp)
	})
}
