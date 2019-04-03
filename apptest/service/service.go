package service

import "context"

type Service interface {
	Serve(ctx context.Context, data ...interface{}) error
}
