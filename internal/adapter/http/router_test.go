package http

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

type stubPingService struct {
	message string
	err     error
}

func (s stubPingService) Ping(ctx context.Context) (string, error) {
	return s.message, s.err
}

func TestHandlePing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := NewRouter(stubPingService{message: "pong"})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/ping", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", resp.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("unexpected error parsing response: %v", err)
	}

	if body["message"] != "pong" {
		t.Fatalf("unexpected message: %v", body["message"])
	}
}

func TestHandlePingServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := NewRouter(stubPingService{err: context.Canceled})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/ping", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusRequestTimeout {
		t.Fatalf("expected timeout status, got %d", resp.Code)
	}
}
