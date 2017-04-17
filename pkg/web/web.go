package web

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/urfave/negroni"

	"github.com/yuuki/hworq/pkg/db"
)

const (
	DefaultShutdownTimeout = 10
)

// Handler serves various HTTP endpoints of the Diamond server
type Handler struct {
	server *http.Server
	db     *db.DB
}

// Option for the web Handler.
type Option struct {
	Port string
	DB   *db.DB
}

// New initializes a new web Handler.
func New(o *Option) *Handler {
	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())

	srv := &http.Server{Addr: ":" + o.Port, Handler: n}

	h := &Handler{
		server: srv,
		db:     o.DB,
	}

	router := httprouter.New()
	router.GET("/ping", h.pingHandler)
	n.UseHandler(router)

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

// pingHandler returns a HTTP handler for the endpoint to ping storage.
func (h *Handler) pingHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := h.db.Ping(); err != nil {
		log.Printf("%+v", err) // Print stack trace by pkg/errors
		unavaliableError(w, errors.Cause(err).Error())
		return
	}
	ok(w, "PONG")
	return
}
