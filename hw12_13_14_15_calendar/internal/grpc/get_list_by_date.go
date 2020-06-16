package grpc

import (
	"context"
	"time"

	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/schema"

	"github.com/pkg/errors"
)

// Get events list by date.
func (s *Server) GetListByDate(ctx context.Context, req *schema.GetListByDateRpcRequest) (*schema.GetListByDateRpcResponse, error) {
	if req.Date == nil {
		return nil, ErrDateNotSet
	}

	events, err := s.r.GetListByDate(ctx, time.Unix(req.Date.Seconds, int64(req.Date.Nanos)))
	if err != nil {
		return nil, errors.Wrap(err, "can't get events list by date")
	}

	resEvents := s.prepareEventsResponse(events)
	return &schema.GetListByDateRpcResponse{
		Event: resEvents,
	}, nil
}
