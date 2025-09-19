package main

import (
	httpAdapter "AgilityFeat-Backend/internal/adapter/http"
	"AgilityFeat-Backend/internal/app/ping"
	"AgilityFeat-Backend/internal/app/underwriting"
	"errors"
	"log"
	"net/http"
)

func main() {
	pingService := ping.NewService()
	underwritingService := underwriting.NewService()

	router := httpAdapter.NewRouter(pingService, underwritingService)

	if err := router.Run("0.0.0.0:8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server exited with error: %v", err)
	}

}
