package http

import (
	"AgilityFeat-Backend/internal/app/underwriting"
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
	result underwriting.Result
	err    error
}

func (s stubPingService) Ping(ctx context.Context) (string, error) {
	return s.message, s.err
}

func (s stubUnderwritingService) Evaluate(ctx context.Context, input underwriting.ApplicationInput) (underwriting.Result, error) {
	return s.result, s.err
}

func TestHandleUnderWriting(t *testing.T) {
	gin.SetMode(gin.TestMode)
	response := underwriting.Result{
		Decision: "Approve",
		LTV:      0.0,
		DTI:      0.0,
		Reasons:  make([]string, 0),
	}
	router := NewRouter(stubPingService{}, stubUnderwritingService{result: response})
	payload := map[string]any{
		"monthly_income": 0,
		"monthly_debts":  0,
		"loan_amount":    0,
		"property_value": 0,
		"credit_score":   0,
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
