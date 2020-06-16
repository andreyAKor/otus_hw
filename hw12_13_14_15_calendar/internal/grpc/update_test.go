package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/repository/repository"
	repositoryMocks "github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/repository/repository/mocks"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/schema"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	ev := repository.Event{
		Title:    "title-123",
		Date:     time.Now().Round(0),
		Duration: time.Duration(time.Hour),
		UserID:   32,
	}

	t.Run("empty", func(t *testing.T) {
		t.Run("event request not set", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := repositoryMocks.NewMockEventsRepo(ctrl)

			srv, err := New(repo, "", 0)
			require.NoError(t, err)

			_, err = srv.Update(context.Background(), &schema.UpdateRpcRequest{})
			require.Equal(t, ErrEventNotSet, err)
		})
	})
	t.Run("errors", func(t *testing.T) {
		t.Run("can't create new event", func(t *testing.T) {
			t.Run("can't create new event", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				repo := repositoryMocks.NewMockEventsRepo(ctrl)

				srv, err := New(repo, "", 0)
				require.NoError(t, err)

				repo.EXPECT().Update(context.Background(), int64(0), ev).Return(repository.ErrTimeBusy)

				resp, err := srv.Update(context.Background(), prepareUpdateRpcRequest(ev))
				require.Equal(t, repository.ErrTimeBusy, errors.Cause(err))
				require.Nil(t, resp)
			})
			t.Run("event not found", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				repo := repositoryMocks.NewMockEventsRepo(ctrl)

				srv, err := New(repo, "", 0)
				require.NoError(t, err)

				repo.EXPECT().Update(context.Background(), int64(0), ev).Return(repository.ErrNotFound)

				resp, err := srv.Update(context.Background(), prepareUpdateRpcRequest(ev))
				require.Equal(t, repository.ErrNotFound, errors.Cause(err))
				require.Nil(t, resp)
			})
			t.Run("other errors", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				repo := repositoryMocks.NewMockEventsRepo(ctrl)

				srv, err := New(repo, "", 0)
				require.NoError(t, err)

				repo.EXPECT().Update(context.Background(), int64(0), ev).Return(errors.New(""))

				resp, err := srv.Update(context.Background(), prepareUpdateRpcRequest(ev))
				require.Error(t, err)
				require.Nil(t, resp)
			})
		})
	})
	t.Run("normal", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositoryMocks.NewMockEventsRepo(ctrl)

		srv, err := New(repo, "", 0)
		require.NoError(t, err)

		repo.EXPECT().Update(context.Background(), int64(0), ev).Return(nil)

		resp, err := srv.Update(context.Background(), prepareUpdateRpcRequest(ev))
		require.NoError(t, err)
		require.EqualValues(t, empty.Empty{}, *resp)
	})
}

func prepareUpdateRpcRequest(ev repository.Event) *schema.UpdateRpcRequest {
	return &schema.UpdateRpcRequest{
		Event: &schema.Event{
			Title: ev.Title,
			Date: &timestamp.Timestamp{
				Seconds: ev.Date.Unix(),
				Nanos:   int32(ev.Date.Nanosecond()),
			},
			Duration: &duration.Duration{
				Seconds: int64(ev.Duration.Seconds()),
				Nanos:   int32(ev.Duration % time.Second),
			},
			UserID: ev.UserID,
		},
	}
}
