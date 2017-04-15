package web

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/urfave/negroni"
)

const (
	DefaultShutdownTimeout = 10
)

// Handler serves various HTTP endpoints of the Diamond server
type Handler struct {
	server *http.Server
}

// Options for the web Handler.
type Option struct {
	Port string
}

// New initializes a new web Handler.
func New(o *Option) *Handler {
	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())

	srv := &http.Server{Addr: ":" + o.Port, Handler: n}

	h := &Handler{
		server: srv,
	}

	mux := http.NewServeMux()
	n.UseHandler(mux)

	return h
}

// Run serves the HTTP endpoints.
func (h *Handler) Run() {
	log.Printf("Listening on :%s\n", h.server.Addr)
	if err := h.server.ListenAndServe(); err != nil {
		log.Println(err)
	}
}

// Shutdown shoudowns the HTTP server.
func (h *Handler) Shutdown(sig os.Signal) error {
	log.Printf("Received %s gracefully shutdown...\n", sig)
	ctx, cancel := context.WithTimeout(context.Background(), DefaultShutdownTimeout)
	defer cancel()
	if err := h.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
