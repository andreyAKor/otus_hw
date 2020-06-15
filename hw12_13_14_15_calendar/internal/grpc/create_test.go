package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/repository"
	repositoryMocks "github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/repository/mocks"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/schema"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
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

			_, err = srv.Create(context.Background(), &schema.CreateRpcRequest{})
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

				repo.EXPECT().Create(context.Background(), ev).Return(int64(0), repository.ErrTimeBusy)

				resp, err := srv.Create(context.Background(), prepareCreateRpcRequest(ev))
				require.Equal(t, repository.ErrTimeBusy, errors.Cause(err))
				require.Nil(t, resp)
			})
			t.Run("other errors", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				repo := repositoryMocks.NewMockEventsRepo(ctrl)

				srv, err := New(repo, "", 0)
				require.NoError(t, err)

				repo.EXPECT().Create(context.Background(), ev).Return(int64(0), errors.New(""))

				resp, err := srv.Create(context.Background(), prepareCreateRpcRequest(ev))
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

		id := int64(1)

		repo.EXPECT().Create(context.Background(), ev).Return(id, nil)

		resp, err := srv.Create(context.Background(), prepareCreateRpcRequest(ev))
		require.NoError(t, err)
		require.EqualValues(t, schema.CreateRpcResponse{
			Id: id,
		}, *resp)
	})
}

func prepareCreateRpcRequest(ev repository.Event) *schema.CreateRpcRequest {
	return &schema.CreateRpcRequest{
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
