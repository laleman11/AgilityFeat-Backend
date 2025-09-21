package http

import (
	"AgilityFeat-Backend/internal/port"
	"bytes"
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

type stubUnderwritingService struct {
	result  port.UnderwritingResponse
	err     error
	history []port.UnderwritingHistoryEntry
}

func (s stubPingService) Ping(ctx context.Context) (string, error) {
	return s.message, s.err
}

func (s stubUnderwritingService) Evaluate(ctx context.Context, input port.UnderwritingRequest) (port.UnderwritingResponse, error) {
	return s.result, s.err
}

func (s stubUnderwritingService) History(ctx context.Context, userID string) ([]port.UnderwritingHistoryEntry, error) {
	return s.history, s.err
}

func TestHandleUnderWriting(t *testing.T) {
	gin.SetMode(gin.TestMode)
	response := port.UnderwritingResponse{
		Decision: "Approve",
		LTV:      0.0,
		DTI:      0.0,
		Reasons:  make([]string, 0),
	}
	router := NewRouter(stubPingService{}, stubUnderwritingService{result: response})
	payload := map[string]any{
		"user_id":        "1234",
		"monthly_income": 5000,
		"monthly_debts":  2000,
		"loan_amount":    270000,
		"property_value": 300000,
		"credit_score":   500,
		"occupancy_type": "primary",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/underwriting", bytes.NewReader(body))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", resp.Code)
	}
	var bodyResp map[string]any
	if err := json.Unmarshal(resp.Body.Bytes(), &bodyResp); err != nil {
		t.Fatalf("unexpected error parsing response: %v", err)
	}

	if bodyResp["decision"] != "Approve" {
		t.Fatalf("unexpected decision: %v", bodyResp["decision"])
	}
}

func TestHandleUnderwritingError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := NewRouter(stubPingService{}, stubUnderwritingService{err: context.Canceled})
	payload := map[string]any{
		"monthly_income": 0,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/underwriting", bytes.NewReader(body))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected timeout status, got %d", resp.Code)
	}
}

func TestHandleUnderwritingHistory(t *testing.T) {
	gin.SetMode(gin.TestMode)

	history := []port.UnderwritingHistoryEntry{
		{
			UserID:   "user-123",
			Response: port.UnderwritingResponse{Decision: "Approve"},
		},
	}

	router := NewRouter(stubPingService{}, stubUnderwritingService{history: history})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/underwriting/history/user-123", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", resp.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("unexpected error parsing response: %v", err)
	}

	items, ok := body["items"].([]any)
	if !ok || len(items) != 1 {
		t.Fatalf("expected one history item, got %v", body["items"])
	}
}

func TestHandleUnderwritingHistoryServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := NewRouter(stubPingService{}, stubUnderwritingService{err: context.Canceled})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/underwriting/history/user-123", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusRequestTimeout {
		t.Fatalf("expected timeout status, got %d", resp.Code)
	}
}

func TestHandlePing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := NewRouter(stubPingService{message: "pong"}, stubUnderwritingService{})

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

	router := NewRouter(stubPingService{err: context.Canceled}, stubUnderwritingService{})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/ping", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusRequestTimeout {
		t.Fatalf("expected timeout status, got %d", resp.Code)
	}
}
