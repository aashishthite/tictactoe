package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

const (
	DEFAULT_PORT = ":8080"
	timeout      = 15 * time.Second
)

type Handlers struct{}

func New() *Handlers {
	return &Handlers{}
}

func (h *Handlers) ListenAndServe() error {
	log.Println("Listening on port", DEFAULT_PORT)
	return http.ListenAndServe(DEFAULT_PORT, h.endpoints())
}

func (h *Handlers) endpoints() *mux.Router {
	router := mux.NewRouter().StrictSlash(false)

	addHandler(router, "GET", "/", h.Health)

	return router
}

// Health check endpoint
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "All OK")
}

//Helper Funcs
func addHandler(r *mux.Router, method, pattern string, hand func(w http.ResponseWriter, r *http.Request)) {
	r.HandleFunc(pattern, requestLogger(hand)).Methods(method)
}

func requestLogger(hand func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("source:%s\tmethod:%s\tpath:%s\n", strings.Split(r.RemoteAddr, ":")[0], r.Method, r.URL.Path)

		done := make(chan bool, 1)
		go func() {
			hand(w, r)
			done <- true
		}()

		select {
		case <-done:
			return
		case <-time.After(timeout - time.Second):
			log.Printf("TIMEOUT ERROR: source:%s\tmethod:%s\tpath:%s", strings.Split(r.RemoteAddr, ":")[0], r.Method, r.URL.Path)
		}
	}
}
