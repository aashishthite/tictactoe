package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aashishthite/tictactoe/pkg/core"
	"github.com/gorilla/mux"
	"github.com/nlopes/slack"
)

const (
	DEFAULT_PORT = "8080"
	timeout      = 15 * time.Second
	SLACK_TOKEN  = ""
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
	port := os.Getenv("PORT")

	if port == "" {
		port = DEFAULT_PORT
	}
	log.Println("Listening on port", port)
	return http.ListenAndServe(":"+port, h.endpoints())
}

func (h *Handlers) endpoints() *mux.Router {
	router := mux.NewRouter().StrictSlash(false)

	addHandler(router, "GET", "/", h.Health)
	/*
		addHandler(router, "POST", "/help", h.Help)

		addHandler(router, "POST", "/start", h.Start)

		addHandler(router, "POST", "/move", h.Move)

		addHandler(router, "POST", "/state", h.State)
	*/
	addHandler(router, "POST", "/cmd", h.cmdHandler)
	//addHandler(router, "POST", "/forfiet", h.Forfiet)

	return router
}

// Health check endpoint
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "All OK")
}

func (h *Handlers) cmdHandler(w http.ResponseWriter, r *http.Request) {
	var v url.Values
	_ = r.ParseForm()
	v = r.Form

	sc := &SlashCommand{
		Token:       v.Get("token"),
		TeamId:      v.Get("team_id"),
		TeamDomain:  v.Get("team_domain"),
		ChannelId:   v.Get("channel_id"),
		ChannelName: v.Get("channel_name"),
		UserId:      v.Get("user_id"),
		UserName:    v.Get("user_name"),
		Command:     v.Get("command"),
		Text:        v.Get("text"),
		ResponseURL: v.Get("response_url"),
	}

	log.Println(sc.Text)
	retval := ""

	switch {
	case strings.ToLower(sc.Text) == "help":
		retval = h.Help()
	case strings.ToLower(sc.Text) == "state":
		retval = h.State()
	default:
		retval = "No Command Found"
	}

	scm := &SlashCommandMessage{
		ResponseType: "in_channel",
		Text:         retval,
		Attachments:  []Attachment{},
	}

	var jsonData bytes.Buffer
	if err := json.NewEncoder(&jsonData).Encode(scm); err != nil {
		return
	}

	cl := &http.Client{}
	resp, err := cl.Post(sc.ResponseURL, "application/json; charset=utf-8", &jsonData)

	if err != nil {
		log.Println(err)
	}
	resp.Body.Close()
}

// Help endpoint
func (h *Handlers) Help() string {
	return "```/ttt @userid : Starts a game against userid" +
		"\n" +
		"/ttt state      : Display state of the current game" +
		"```"
}

func (h *Handlers) Start(user1, user2 string) string {

	api := slack.New(SLACK_TOKEN)
	//check if valid user2
	valid := false
	users, err := api.GetUsers()
	if err != nil {
		return "unable to validate user"
	}
	for _, v := range users {
		if user2 == v.Name {
			valid = true
		}
	}
	if !valid {
		return "Not a valid user"
	}
	return "Starting a game:"
	/*
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
	*/
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

func (h *Handlers) State() string {

	s := h.GameEngine.GameState()
	return "```" + s.GameBoard + "```" + "\n" + s.Status
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
