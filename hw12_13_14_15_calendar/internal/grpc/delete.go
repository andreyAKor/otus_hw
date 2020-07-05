package grpc

import (
	"context"

	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/schema"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
)

// Delete event.
func (s *Server) Delete(ctx context.Context, req *schema.DeleteRpcRequest) (*empty.Empty, error) {
	if err := s.calendar.Delete(ctx, req.Id); err != nil {
		return nil, errors.Wrap(err, "can't delete event")
	}

	return &empty.Empty{}, nil
}
