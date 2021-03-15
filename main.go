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
	if arg == "build" {
		fmt.Println("Building...")
		if err := ssg.Build(); err != nil {
			log.Fatal(err)
		}
		return
	} else if arg == "dev" || arg == "start" {
		ctx := context.Background()
		var h http.Handler
		var err error
		if arg == "dev" {
			fmt.Println("Running dev server...")
			h, err = ssr.New(ctx)
			if err != nil {
				log.Fatal(err)
				return
			}
		} else {
			fmt.Println("Serving files...")
			h, err = serve.New()
			if err != nil {
				log.Fatal(err)
				return
			}
		}
		port := os.Getenv("PORT")
		if port == "" {
			port = "3000"
		}
		srv := &http.Server{
			Handler:      h,
			Addr:         "127.0.0.1:" + port,
			ReadTimeout:  20 * time.Second,
			WriteTimeout: 20 * time.Second,
		}
		go gracefullyShutDownServerOnSignal(srv, ctx)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Unable to to start server: %v", err)
		}
		return
	}
	printHelp()
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

func gracefullyShutDownServerOnSignal(server *http.Server, ctx context.Context) {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGINT)
	<-exit
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Unable to shut down server: %v", err)
	}
}
