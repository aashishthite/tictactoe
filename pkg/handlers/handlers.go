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

	addHandler(router, "POST", "/cmd", h.cmdHandler)

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
		break
	case strings.ToLower(sc.Text) == "state":
		retval = h.State(sc.ChannelId)
		break
	case sc.Text[0] == '@':
		retval = h.Start(sc.UserName, sc.Text[1:len(sc.Text)], sc.ChannelId)
		break
	case len(sc.Text) > 5 && strings.ToLower(sc.Text[0:4]) == "move":
		retval = h.Move(sc.UserName, []rune(sc.Text)[5], sc.ChannelId)
		break
	case strings.ToLower(sc.Text) == "forfeit":
		retval = h.GameEngine.Forfeit(sc.UserName, sc.ChannelId)
		break
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
	return "```/ttt @userid : Starts a game against userid\n" +
		"/ttt state      : Display state of the current game\n" +
		"/ttt move <pos> : Makes a valid move to pos\n" +
		"/ttt forfeit    : Forfeits the current game\n" +
		"```"
}

func (h *Handlers) Start(user1, user2, channelId string) string {
	if user1 == user2 {
		return "Go find friends to play tictactoe with"
	}

	api := slack.New(os.Getenv("SLACK_API_TOKEN"))
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

	state, err := h.GameEngine.Start(&core.Player{ID: user1}, &core.Player{ID: user2}, channelId)
	if err != nil {
		return "A game is still going on"
	}

	return "Starting a game: \n" + "```" + state.GameBoard + "```" + "\n" + state.Status

}

func (h *Handlers) Move(user string, pos rune, channelId string) string {
	state, err := h.GameEngine.Move(&core.Player{ID: user}, pos, channelId)
	if err != nil {
		return err.Error()
	}

	return "```" + state.GameBoard + "```" + "\n" + state.Status
}

func (h *Handlers) State(channelId string) string {

	s := h.GameEngine.GameState(channelId)
	return "```" + s.GameBoard + "```" + "\n" + s.Status
}

func (h *Handlers) Forfiet(player, channelId string) string {
	return h.GameEngine.Forfeit(player, channelId)
}

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
