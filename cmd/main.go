package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	api "github.com/adityadafe/kc-backend-assgn/internal/api/handlers"
	"github.com/adityadafe/kc-backend-assgn/internal/storage"
)

var bindAddress = ":9090"

func main() {

	//for now just print logs in fd 1
	logger := log.New(os.Stdout, "image-api ", log.LstdFlags)

	sm := http.NewServeMux()
	store := storage.CreateNewStore()

	getJobInfoHandler := api.NewGetJobInfoHandler(logger, &store)
	submitJobHandler := api.NewSubmitJobHandler(logger, &store)

	sm.Handle("POST /api/submit", submitJobHandler)
	sm.Handle("GET /api/status", getJobInfoHandler)

	server := &http.Server{
		Addr:         bindAddress,
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		logger.Println("Starting server in another routine")
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	sig := <-sigChan

	logger.Println("Received termination, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 5*time.Second)
	server.Shutdown(tc)

}
