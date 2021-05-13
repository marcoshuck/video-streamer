package main

import (
	"github.com/marcoshuck/video-streamer/pkg/server"
	"log"
)

func main() {
	log.SetPrefix("[APP] ")
	log.Println("Initializing application")
	if err := run(); err != nil {
		log.Fatalf("error while running application: %s", err)
	}
}

func run() error {
	s := server.NewServer()
	log.Println("List of available routes:")
	if err := s.WalkRoutes(); err != nil {
		return err
	}
	log.Println("Running HTTP server on port 8080")
	return s.ListenAndServe(8080)
}
