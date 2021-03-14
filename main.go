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
			signal.NotifyContext(ctx, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGINT)
		case "start":
			log.Fatal(serve.Serve())
		case "help":
			fallthrough
		default:
			printHelp()
		}
	} else {
		printHelp()
	}
}

func printHelp() {
	help := `
JonahGold's Henkernieuws

Available Commands

start - Serves the built static HTML pages
build - Builds static HTML pages for serving
dev   - Creates a server that dynamically parses pages on request (and runs snowpack)
help  - Prints this help message

`
	fmt.Println(help)
}
