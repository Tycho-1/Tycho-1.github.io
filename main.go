// Command portfolio generates a static GitHub Pages site and optionally serves it locally.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Tycho-1/gh-pages-portfolio/internal/site"
)

func main() {
	out := flag.String("o", "dist", "output directory for generated HTML")
	serve := flag.Bool("serve", false, "serve the generated site after build")
	port := flag.String("port", "8080", "local preview port when -serve is set")
	flag.Parse()

	root, err := site.Root()
	if err != nil {
		log.Fatalf("resolve root: %v", err)
	}

	gen := site.NewGenerator(root, *out)
	if err := gen.Build(); err != nil {
		log.Fatalf("build site: %v", err)
	}

	fmt.Printf("Generated site in %s/\n", *out)

	if !*serve {
		fmt.Printf("Preview locally: go run . -serve\n")
		return
	}

	addr := ":" + *port
	handler := http.FileServer(http.Dir(*out))
	server := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	fmt.Printf("Serving at http://localhost%s/\n", addr)
	fmt.Println("Press Ctrl+C to stop.")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("serve site: %v", err)
	}
}

func init() {
	log.SetOutput(os.Stderr)
	log.SetFlags(0)
}
