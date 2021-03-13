package main

import (
	"fmt"
	"log"
	"os"

	"github.com/theonejonahgold/pwa/serve"
	"github.com/theonejonahgold/pwa/ssg"
	"github.com/theonejonahgold/pwa/ssr"
)

func main() {
	if len(os.Args) > 1 {
		arg := os.Args[1]
		switch arg {
		case "ssg":
			if err := ssg.SSG(); err != nil {
				log.Fatal(err)
			}
		case "ssr":
			log.Fatal(ssr.SSR())
		case "serve":
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

ssg   - Builds static HTML pages for serving
ssr   - Creates a server that dynamically parses pages on request
serve - Serves the built static HTML pages
help  - Prints this help message\n
	`
	fmt.Print(help)
}
