package ping

import "context"

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Ping(ctx context.Context) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	return "pong", nil
}
