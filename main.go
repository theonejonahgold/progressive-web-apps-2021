package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/theonejonahgold/pwa/serve"
	"github.com/theonejonahgold/pwa/ssg"
	"github.com/theonejonahgold/pwa/ssr"
)

func main() {
	fmt.Println(`
Jonahgold's Henkernieuws`)
	if len(os.Args) > 1 {
		arg := os.Args[1]
		switch arg {
		case "build":
			if err := ssg.SSG(); err != nil {
				log.Fatal(err)
			}
		case "dev":
			ctx, err := ssr.SSR()
			if err != nil {
				log.Fatal(err)
			}
			_, stop := signal.NotifyContext(ctx, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGINT)
			defer stop()
		case "start":
			if err := serve.Serve(); err != nil {
				fmt.Println(err)
			}
		default:
			printHelp()
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
