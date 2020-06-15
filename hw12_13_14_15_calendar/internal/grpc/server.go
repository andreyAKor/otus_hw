package grpc

import (
	"context"
	"net"
	"strconv"
	"time"

	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/repository"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/schema"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

var (
	ErrDateNotSet  = errors.New("date not set")
	ErrStartNotSet = errors.New("start not set")
	ErrEventNotSet = errors.New("event request not set")
)

var _ schema.CalendarServer = (*Server)(nil)

//go:generate protoc --proto_path=../../schema --go_out=plugins=grpc:../../schema ../../schema/calendar.proto
type Server struct {
	r    repository.EventsRepo
	host string
	port int
}

func New(r repository.EventsRepo, host string, port int) (*Server, error) {
	return &Server{r, host, port}, nil
}

// Running grpc-server.
func (s *Server) Run(ctx context.Context) error {
	lis, err := net.Listen("tcp", net.JoinHostPort(s.host, strconv.Itoa(s.port)))
	if err != nil {
		return errors.Wrap(err, "grpc-server listen fail")
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(s.unaryInterceptor),
	)

	schema.RegisterCalendarServer(grpcServer, s)
	if err := grpcServer.Serve(lis); err != nil {
		return errors.Wrap(err, "grpc-server serve fail")
	}

	return nil
}

// Logging unary interceptor function to handle logging per RPC call.
func (s Server) unaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	// Calls the handler
	h, err := handler(ctx, req)

	i := log.Info()

	if peer, ok := peer.FromContext(ctx); ok {
		host, _, err := net.SplitHostPort(peer.Addr.String())
		if err != nil {
			return h, err
		}

		i.Str("ip", host)
	}

	i.Str("startAt", start.String()).
		Str("method", info.FullMethod).
		Interface("request", req).
		TimeDiff("latency", time.Now(), start).
		Err(err).
		Msg("grpc-request")

	return h, err
}
