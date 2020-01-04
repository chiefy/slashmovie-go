package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/chiefy/go-slack-utils/pkg/blockui"
	tmdb "github.com/ryanbradynd05/go-tmdb"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	// TmdbAPIKey is the env var name containing the TMDB API key
	TmdbAPIKey       = "TMDB_API_KEY"
	tmdbImageURLBase = "https://image.tmdb.org/t/p"
	tmdbImageSmall   = "w92"
	tmdbImageLarge   = "w342"
	numResults       = 5
)

var (
	tmdbAPI *tmdb.TMDb
)

func init() {
	k := os.Getenv(TmdbAPIKey)
	if k == "" {
		log.Fatalf("no env var found for %s", TmdbAPIKey)
	}
	tmdbAPI = tmdb.Init(k)
}

func MovieLookupHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var action *SlackBlockAction
	err := json.Unmarshal([]byte(r.Form.Get("payload")), &action)
	if err != nil {
		log.Printf("error decoding JSON from payload - %s", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	id, _ := strconv.Atoi(action.Actions[0].Value)
	movie, err := tmdbAPI.GetMovieInfo(id, map[string]string{})
	if err != nil {
		log.Printf("error looking up movie ID - %s", err)
		http.Error(w, "Bad Movie Lookup", http.StatusInternalServerError)
		return
	}
	sm := NewSlackMessage()
	sm.ResponseType = "in_channel"
	tb := &SlackBlock{
		Type: "section",
		Text: &SlackText{
			Type: "mrkdwn",
			Text: makeMovieMarkdown(movie),
		},
		Accessory: &SlackAccessory{
			Type:     "image",
			ImageURL: fmt.Sprintf("%s/%s/%s", tmdbImageURLBase, tmdbImageSmall, movie.PosterPath),
			AltText:  movie.Title,
		},
	}

	ib := &SlackBlock{
		Type: "section",
		Accessory: &SlackAccessory{
			Type: "image",
			Text: &SlackText{
				Type:  "plain_text",
				Text:  "",
				Emoji: false,
			},
			ImageURL: fmt.Sprintf("%s/%s/%s", tmdbImageURLBase, tmdbImageLarge, movie.PosterPath),
		},
	}

	sm.AddBlock(tb)
	sm.AddBlock(ib)

	j, err := json.Marshal(sm)
	if err != nil {
		log.Printf("error marshalling json %s", err)
		http.Error(w, "JSON Error", http.StatusInternalServerError)
		return
	}
	url := action.ResponseURL
	c := &http.Client{}
	log.Println(url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", action.Token))

	log.Println(string(j))
	res, err := c.Do(req)
	if err != nil {
		log.Printf("error sending response %s", err)
		http.Error(w, "Bad Request", http.StatusInternalServerError)
		return
	}
	log.Printf("%#v", res)

}

// MovieSearchHandler handles the slack slash command POST request
func MovieSearchHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println("Got request to moviehandler")
	movieStr := r.FormValue("text")
	opts := map[string]string{}

	res, err := tmdbAPI.SearchMovie(movieStr, opts)
	if err != nil {
		log.Println(err)
		http.Error(w, "movie lookup error", http.StatusInternalServerError)
		return
	}
	sm := NewSlackMessage()
	tb := &SlackBlock{
		Type: "section",
		Text: &SlackText{
			Type: "mrkdwn",
			Text: fmt.Sprintf("Found %d results for *\"%s\"*:", len(res.Results), movieStr),
		},
	}

	sm.AddBlock(tb)
	sm.AddDivider()
	buttons := make([]SlackAccessory, 0)

	for i, m := range res.Results {
		y := strings.Split(m.ReleaseDate, "-")
		button := &SlackAccessory{
			Type: "button",
			Text: &SlackText{
				Type: "plain_text",
				Text: fmt.Sprintf("%s (%s)", m.Title, y[0]),
			},
			Value: strconv.Itoa(m.ID),
		}
		buttons = append(buttons, *button)
		if i >= numResults-1 {
			break
		}
	}

	sm.AddBlock(&SlackBlock{
		Type:     "actions",
		Elements: buttons,
	})

	j, err := json.Marshal(sm)
	log.Println(string(j))
	if err != nil {
		log.Printf("error marshalling json %s", err)
		http.Error(w, "JSON Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)

}

func makeMovieMarkdown(movie *tmdb.Movie) string {
	l := fmt.Sprintf("<https://www.imdb.com/title/%s|IMDB link>", movie.ImdbID)
	return fmt.Sprintf("*%s*\nRelease Date: %s\nBudget: $%d\n>%s\n\n%s",
		movie.Title, movie.ReleaseDate, movie.Budget, movie.Overview, l)
}
