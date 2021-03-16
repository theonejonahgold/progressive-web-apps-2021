package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// RunContext runs a http server with the provided handler and context.
//
// The context's channel will be closed when the server is shut down.
func RunContext(h http.Handler, ctx context.Context) {
	srv := createServer(h)
	go gracefullyShutDownOnSignal(srv, ctx)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Unable to to start server: %v", err)
	}
}

// Run runs a http with the provided handler
func Run(h http.Handler) {
	srv := createServer(h)
	go gracefullyShutDownOnSignal(srv, context.Background())
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Unable to to start server: %v", err)
	}
}

func createServer(h http.Handler) *http.Server {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("Starting server on port :%v", port)
	srv := &http.Server{
		Handler:      h,
		Addr:         ":" + port,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}
	return srv
}

func gracefullyShutDownOnSignal(srv *http.Server, ctx context.Context) {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGINT)
	sigcall := <-exit
	log.Printf("%v - shutting down server...", sigcall)
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Unable to shut down server: %v", err)
	}
}
