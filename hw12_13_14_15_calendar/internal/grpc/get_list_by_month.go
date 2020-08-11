package grpc

import (
	"context"
	"time"

	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/schema"

	"github.com/pkg/errors"
)

// Get events list by month.
func (s *Server) GetListByMonth(ctx context.Context, req *schema.GetListByMonthRpcRequest) (*schema.GetListByMonthRpcResponse, error) {
	if req.Start == nil {
		return nil, ErrStartNotSet
	}

	events, err := s.calendar.GetListByMonth(ctx, time.Unix(req.Start.Seconds, int64(req.Start.Nanos)))
	if err != nil {
		return nil, errors.Wrap(err, "can't get events list by month")
	}

	resEvents := s.prepareEventsResponse(events)
	return &schema.GetListByMonthRpcResponse{
		Event: resEvents,
	}, nil
}
