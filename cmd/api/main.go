package main

import (
	httpAdapter "AgilityFeat-Backend/internal/adapter/http"
	"AgilityFeat-Backend/internal/app/ping"
	"AgilityFeat-Backend/internal/app/underwriting"
	underwritingRepository "AgilityFeat-Backend/internal/infra/memory/underwriting"
	"errors"
	"log"
	"net/http"
)

func main() {
	pingService := ping.NewService()
	repository := underwritingRepository.NewRepository()
	underwritingService := underwriting.NewService(repository)

	router := httpAdapter.NewRouter(pingService, underwritingService)

	if err := router.Run("0.0.0.0:8081"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server exited with error: %v", err)
	}
}
