package http

import (
	underwritingapp "AgilityFeat-Backend/internal/app/underwriting"
	"AgilityFeat-Backend/internal/port"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type handler struct {
	pingService         port.PingService
	underwritingService port.UnderWritingService
}

func NewRouter(pingService port.PingService, underwritingService port.UnderWritingService) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	h := handler{pingService, underwritingService}

	v1 := router.Group("/api/v1")
	{
		v1.GET("/ping", h.handlePing)
		v1.POST("underwriting", h.handleUnderwriting)
	}

	return router
}

type underwritingRequest struct {
	MonthlyIncome float64 `json:"monthly_income" binding:"required"`
	MonthlyDebts  float64 `json:"monthly_debts" binding:"required" `
	LoanAmount    float64 `json:"loan_amount" binding:"required"`
	PropertyValue float64 `json:"property_value" binding:"required"`
	CreditScore   int     `json:"credit_score" binding:"required"`
	OccupancyType string  `json:"occupancy_type" binding:"required"`
}

func (h handler) handleUnderwriting(c *gin.Context) {
	var payload underwritingRequest
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	result, err := h.underwritingService.Evaluate(c.Request.Context(), underwritingapp.ApplicationInput{
		MonthlyIncome: payload.MonthlyIncome,
		MonthlyDebts:  payload.MonthlyDebts,
		LoanAmount:    payload.LoanAmount,
		PropertyValue: payload.PropertyValue,
		CreditScore:   payload.CreditScore,
		OccupancyType: payload.OccupancyType,
	})

	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			status = http.StatusRequestTimeout
		}
		c.JSON(status, gin.H{"error": http.StatusText(status)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"decision": result.Decision,
		"dti":      result.DTI,
		"ltv":      result.LTV,
		"reasons":  result.Reasons,
	})
}

func (h handler) handlePing(c *gin.Context) {
	if h.pingService == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": http.StatusText(http.StatusNotFound)})
		return
	}

	message, err := h.pingService.Ping(c.Request.Context())
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			status = http.StatusRequestTimeout
		}
		c.JSON(status, gin.H{"error": http.StatusText(status)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": message})
}
