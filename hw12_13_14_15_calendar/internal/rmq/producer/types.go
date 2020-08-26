package producer

import (
	"context"
	"io"
)

type Producer interface {
	Run(ctx context.Context) error
	Publish(content []byte) error
	io.Closer
}
