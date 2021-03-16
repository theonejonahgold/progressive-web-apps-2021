package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/theonejonahgold/pwa/prerender"
	"github.com/theonejonahgold/pwa/server"
	"github.com/theonejonahgold/pwa/server/ssr"
	"github.com/theonejonahgold/pwa/server/static"
)

func main() {
	if len(os.Args) <= 1 {
		printHelp()
		return
	}
	arg := os.Args[1]
	switch arg {
	case "build":
		if err := prerender.Build(); err != nil {
			log.Fatal(err)
		}
	case "dev":
		ctx := context.Background()
		server.RunContext(ssr.New(ctx), ctx)
	case "start":
		server.Run(static.New())
	case "build-cron":
		if err := prerender.ScheduleBuild(); err != nil {
			log.Fatal(err)
		}
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
