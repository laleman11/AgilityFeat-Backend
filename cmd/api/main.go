package main

import (
	httpAdapter "AgilityFeat-Backend/internal/adapter/http"
	"errors"
	"log"
	"net/http"
)

func main() {

	router := httpAdapter.NewRouter()

	if err := router.Run("0.0.0.0:8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server exited with error: %v", err)
	}

}
