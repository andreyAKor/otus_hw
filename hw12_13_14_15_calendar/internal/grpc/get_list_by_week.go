package grpc

import (
	"context"
	"time"

	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/schema"

	"github.com/pkg/errors"
)

// Get events list by week.
func (s *Server) GetListByWeek(ctx context.Context, req *schema.GetListByWeekRpcRequest) (*schema.GetListByWeekRpcResponse, error) {
	if req.Start == nil {
		return nil, ErrStartNotSet
	}

	events, err := s.calendar.GetListByWeek(ctx, time.Unix(req.Start.Seconds, int64(req.Start.Nanos)))
	if err != nil {
		return nil, errors.Wrap(err, "can't get events list by week")
	}

	resEvents := s.prepareEventsResponse(events)
	return &schema.GetListByWeekRpcResponse{
		Event: resEvents,
	}, nil
}
