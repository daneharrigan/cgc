package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/daneharrigan/cgc/counter"
)

var (
	BungieAPIKey   = os.Getenv("BUNGIE_API_KEY")
	ErrInvalidFrom = errors.New("invalid from parameter")
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("ns=cgc ")
}

func main() {
	log.Print("at=start")
	addr := fmt.Sprintf(":%s", GetenvOrDefault("PORT", "5000"))
	if err := http.ListenAndServe(addr, NewAPI()); err != nil {
		log.Printf("at=finish error=%q", err)
	}
	log.Print("at=finish")
}

func GetenvOrDefault(key, value string) string {
	if v := os.Getenv("PORT"); v != "" {
		return v
	}

	return value
}

func NewAPI() *API {
	return &API{
		ServeMux: http.NewServeMux(),
		pattern:  regexp.MustCompile(`/games/(\d+)/(\d+)/(\d+)$`),
	}
}

type Message struct {
	message string `json:"message"`
}

type API struct {
	*http.ServeMux
	pattern *regexp.Regexp
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer api.HandleRecover(w, r)

	if api.pattern.MatchString(r.URL.Path) {
		api.HandleGames(w, r)
		return
	}

	api.HandleMessage(404, http.StatusText(404), w, r)
}

func (api *API) HandleMessage(statusCode int, message string, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")

	msg := &Message{message}
	err := json.NewEncoder(w).Encode(msg)
	if err != nil {
		log.Printf("at=HandleMessage fn=Encode error=%q", err)
	}
}

func (api *API) HandleGames(w http.ResponseWriter, r *http.Request) {
	fromParam := r.FormValue("from")

	to := time.Now()
	var from time.Time

	if fromParam == "" {
		from = to
	} else {
		var err error
		from, err = time.Parse(counter.TimeFormat, fromParam)
		if err != nil {
			api.HandleMessage(400, ErrInvalidFrom.Error(), w, r)
			return
		}
	}

	params := api.pattern.FindStringSubmatch(r.URL.Path)
	counter := counter.New(BungieAPIKey, params[1], params[2], params[3])
	results, err := counter.GetResults(from, to)
	if err != nil {
		api.HandleMessage(400, err.Error(), w, r)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		log.Printf("at=HandleGames fn=Encode error=%q", err)
	}
}

func (api *API) HandleRecover(w http.ResponseWriter, r *http.Request) {
	if err := recover(); err != nil {
		log.Printf("at=HandleRecover error=%q", err)
		api.HandleMessage(500, http.StatusText(500), w, r)
	}
}
