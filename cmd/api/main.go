package main

import (
	httpAdapter "AgilityFeat-Backend/internal/adapter/http"
	"AgilityFeat-Backend/internal/app/ping"
	"errors"
	"log"
	"net/http"
)

func main() {
	pingService := ping.NewService()

	router := httpAdapter.NewRouter(pingService)

	if err := router.Run("0.0.0.0:8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server exited with error: %v", err)
	}

}
