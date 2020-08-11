package producer

import (
	"context"
	"io"
)

type Producerer interface {
	Run(ctx context.Context) error
	Publish(content []byte) error
	io.Closer
}
