package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aashishthite/tictactoe/pkg/core"
	"github.com/gorilla/mux"
)

const (
	DEFAULT_PORT = ":8080"
	timeout      = 15 * time.Second
)

type Handlers struct {
	GameEngine *core.Engine
}

func New() *Handlers {
	return &Handlers{
		GameEngine: core.NewEngine(),
	}
}

func (h *Handlers) ListenAndServe() error {
	log.Println("Listening on port", DEFAULT_PORT)
	return http.ListenAndServe(DEFAULT_PORT, h.endpoints())
}

func (h *Handlers) endpoints() *mux.Router {
	router := mux.NewRouter().StrictSlash(false)

	addHandler(router, "GET", "/", h.Health)

	addHandler(router, "GET", "/help", h.Help)

	addHandler(router, "POST", "/start", h.Start)

	addHandler(router, "POST", "/move", h.Move)

	addHandler(router, "GET", "/state", h.State)

	//addHandler(router, "POST", "/forfiet", h.Forfiet)

	return router
}

// Health check endpoint
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "All OK")
}

// Help endpoint
func (h *Handlers) Help(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "/ttt start @userid")
}

func (h *Handlers) Start(w http.ResponseWriter, r *http.Request) {
	var postReq StartGameRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&postReq)
	if err != nil {
		ErrorRespondAndLog(w, r, err, "Failed to decode the request body", http.StatusBadRequest)
		return
	}

	state, err := h.GameEngine.Start(&core.Player{ID: postReq.Player1ID}, &core.Player{ID: postReq.Player1ID})
	if err != nil {
		ErrorRespondAndLog(w, r, err, "A game is still going on", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(state)
}

func (h *Handlers) Move(w http.ResponseWriter, r *http.Request) {
	var postReq MoveRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&postReq)
	if err != nil {
		ErrorRespondAndLog(w, r, err, "Failed to decode the request body", http.StatusBadRequest)
		return
	}

	state, err := h.GameEngine.Move(&core.Player{ID: postReq.PlayerID}, []rune(postReq.Position)[0])
	if err != nil {
		ErrorRespondAndLog(w, r, err, "Invalid Move", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(state)
}

func (h *Handlers) State(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(h.GameEngine.GameState())
}

/*
func (h *Handlers) Forfiet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "All OK")
}
*/

//Helper Funcs

func ErrorRespondAndLog(w http.ResponseWriter, r *http.Request, err error, errStr string, code int) {
	log.Printf("ERROR at path=\"%s\" err=\"%s\" code:\"%d\"", r.RequestURI, err, code)
	http.Error(w, errStr, code)
}

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
