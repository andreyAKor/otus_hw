package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/repository"
	repositoryMocks "github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/repository/mocks"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/schema"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestGetListByDate(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		t.Run("date not set", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := repositoryMocks.NewMockEventsRepo(ctrl)

			srv, err := New(repo, "", 0)
			require.NoError(t, err)

			_, err = srv.GetListByDate(context.Background(), &schema.GetListByDateRpcRequest{})
			require.Equal(t, ErrDateNotSet, err)
		})
	})
	t.Run("errors", func(t *testing.T) {
		t.Run("can't get events list by date", func(t *testing.T) {
			t.Run("other errors", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				repo := repositoryMocks.NewMockEventsRepo(ctrl)

				srv, err := New(repo, "", 0)
				require.NoError(t, err)

				tm := time.Now().Round(0)
				repo.EXPECT().
					GetListByDate(context.Background(), tm).
					Return([]repository.Event{}, errors.New(""))

				resp, err := srv.GetListByDate(context.Background(), &schema.GetListByDateRpcRequest{
					Date: &timestamp.Timestamp{
						Seconds: tm.Unix(),
						Nanos:   int32(tm.Nanosecond()),
					},
				})
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

		tm := time.Now().Round(0)
		repo.EXPECT().
			GetListByDate(context.Background(), tm).
			Return([]repository.Event{}, nil)

		resp, err := srv.GetListByDate(context.Background(), &schema.GetListByDateRpcRequest{
			Date: &timestamp.Timestamp{
				Seconds: tm.Unix(),
				Nanos:   int32(tm.Nanosecond()),
			},
		})
		require.NoError(t, err)
		require.EqualValues(t, schema.GetListByDateRpcResponse{
			Event: []*schema.Event{},
		}, *resp)
	})
}
