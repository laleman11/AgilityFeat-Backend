package ping

import (
	"context"
	"testing"
)

func TestServicePing(t *testing.T) {
	pingService := NewService()
	msg, err := pingService.Ping(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg != "pong" {
		t.Fatalf("unexpected response: %s", msg)
	}
}

func TestServicePingContextCanceled(t *testing.T) {
	pingService := NewService()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := pingService.Ping(ctx); err == nil {
		t.Fatal("expected context error")
	}
}
