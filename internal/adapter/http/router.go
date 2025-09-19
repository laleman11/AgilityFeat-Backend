package http

import (
	"AgilityFeat-Backend/internal/port"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type handler struct {
	pingService port.PingService
}

func NewRouter(pingService port.PingService) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	h := handler{pingService}

	v1 := router.Group("/api/v1")
	{
		v1.GET("/ping", h.handlePing)
	}

	return router
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
