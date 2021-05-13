package server

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	minPort = 49152
	maxPort = 65535
)

type VideoRoutes interface {
	GetVideos(w http.ResponseWriter, r *http.Request)
	StreamVideo(w http.ResponseWriter, r *http.Request)
}

type Server interface {
	ListenAndServe(port uint16) error
	VideoRoutes
	setRouter(router chi.Router) Server
	WalkRoutes() error
}

type server struct {
	router chi.Router
}

func (s server) StreamVideo(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "slug")
	f, err := os.Open(fmt.Sprintf("data/%s.mp4", filename))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to find video: %s", filename), http.StatusNotFound)
		return
	}

	_, err = io.Copy(w, f)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to stream video: %s", filename), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s server) GetVideos(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal([]string{"terrain"})
	if err != nil {
		http.Error(w, "failed to marshal list of videos", http.StatusInternalServerError)
		return
	}

	if n, err := w.Write(data); len(data) != n || err != nil {
		http.Error(w, "failed to send response", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (s *server) setRouter(router chi.Router) Server {
	s.router = router
	return s
}

func (s server) ListenAndServe(port uint16) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", port), s.router)
}

func (s server) WalkRoutes() error {
	return chi.Walk(s.router, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("Method: %s | Route: %s", method, route)
		return nil
	})
}

func NewServer() Server {
	r := newRouter()
	s := newServer()
	r = configureRoutes(s, r)
	s = s.setRouter(r)
	return s
}

func newServer() Server {
	return &server{}
}

func configureRoutes(s VideoRoutes, r chi.Router) chi.Router {
	r.Get("/videos", s.GetVideos)
	r.Get("/videos/{slug}/watch", s.StreamVideo)
	return r
}

func newRouter() chi.Router {
	return chi.NewRouter()
}
