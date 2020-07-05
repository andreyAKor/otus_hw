package grpc

import (
	"context"

	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/schema"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
)

// Updating event.
func (s *Server) Update(ctx context.Context, req *schema.UpdateRpcRequest) (*empty.Empty, error) {
	if req.Event == nil {
		return nil, ErrEventNotSet
	}

	ev := s.prepareEventRequest(req.Event)
	if err := s.calendar.Update(ctx, req.Id, ev); err != nil {
		return nil, errors.Wrap(err, "can't update event")
	}

	return &empty.Empty{}, nil
}
