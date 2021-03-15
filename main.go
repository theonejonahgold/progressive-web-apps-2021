package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/theonejonahgold/pwa/serve"
	"github.com/theonejonahgold/pwa/ssg"
	"github.com/theonejonahgold/pwa/ssr"
)

func main() {
	if len(os.Args) <= 1 {
		printHelp()
		return
	}
	arg := os.Args[1]
	switch arg {
	case "build":
		log.Println("Building...")
		if err := ssg.Build(); err != nil {
			log.Fatal(err)
		}
	case "dev":
		ctx := context.Background()
		log.Println("Running dev server...")
		h, err := ssr.New(ctx)
		if err != nil {
			log.Fatal(err)
			return
		}
		startServer(h, ctx)
	case "start":
		ctx := context.Background()
		log.Println("Serving files...")
		h, err := serve.New()
		if err != nil {
			log.Fatal(err)
			return
		}
		startServer(h, ctx)
	case "build-cron":
		if err := ssg.ScheduleBuild(); err != nil {
			log.Fatal(err)
		}
	}
}

func startServer(h http.Handler, ctx context.Context) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	fmt.Println(port)
	srv := &http.Server{
		Handler:      h,
		Addr:         ":" + port,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}
	go gracefullyShutDownServerOnSignal(srv, ctx)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Unable to to start server: %v", err)
	}
}

func gracefullyShutDownServerOnSignal(server *http.Server, ctx context.Context) {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGINT)
	<-exit
	fmt.Println("Shutting down...")
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Unable to shut down server: %v", err)
	}
}

func printHelp() {
	help := `
Available Commands

start - Serves the built static HTML pages
build - Builds static HTML pages for serving
dev   - Creates a server that dynamically parses pages on request (and runs snowpack)
help  - Prints this help message
`
	fmt.Println(help)
}
