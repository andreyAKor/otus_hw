package grpc

import (
	"context"

	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/schema"

	"github.com/pkg/errors"
)

// Add new event.
func (s *Server) Create(ctx context.Context, req *schema.CreateRpcRequest) (*schema.CreateRpcResponse, error) {
	if req.Event == nil {
		return nil, ErrEventNotSet
	}

	ev := s.prepareEventRequest(req.Event)
	id, err := s.r.Create(ctx, ev)
	if err != nil {
		return nil, errors.Wrap(err, "can't create new event")
	}

	return &schema.CreateRpcResponse{
		Id: id,
	}, nil
}
