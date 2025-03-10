package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	api "github.com/adityadafe/kc-backend-assgn/internal/api/handlers"
)

var bindAddress = ":9090"

func main() {

	logger := log.New(os.Stdout, "image-api ", log.LstdFlags)

	sm := http.NewServeMux()

	getJobInfoHandler := api.NewGetJobInfoHandler(logger)
	SubmitJobHandler := api.NewSubmitJobHandler(logger)

	sm.Handle("POST /api/submit", SubmitJobHandler)
	sm.Handle("GET /api/status", getJobInfoHandler)

	server := &http.Server{
		Addr:         bindAddress,
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		logger.Println("starting server in another routine")
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
