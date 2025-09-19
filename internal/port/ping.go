package port

import "context"

type PingService interface {
	Ping(ctx context.Context) (string, error)
}
