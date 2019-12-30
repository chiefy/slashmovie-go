package main

import (
	"encoding/json"
	"fmt"
	tmdb "github.com/ryanbradynd05/go-tmdb"
	"log"
	"net/http"
	"os"
)

const (
	// TmdbAPIKey is the env var name containing the TMDB API key
	TmdbAPIKey       = "TMDB_API_KEY"
	tmdbImageURLBase = "https://image.tmdb.org/t/p/w92/"
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

// MovieHandler handles the slack slash command POST request
func MovieHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
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

	for i, m := range res.Results {
		t := &SlackText{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*%s* (%s)", m.Title, m.ReleaseDate),
		}
		a := &SlackAccessory{
			Type:     "image",
			ImageURL: fmt.Sprintf("%s/%s", tmdbImageURLBase, m.PosterPath),
			AltText:  m.Title,
		}
		b := &SlackBlock{
			Type:      "section",
			Text:      t,
			Accessory: a,
		}
		sm.AddBlock(b)
		if i >= numResults-1 {
			break
		}
	}

	j, err := json.Marshal(sm)
	if err != nil {
		log.Printf("error marshalling json %s", err)
		http.Error(w, "JSON Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)

}
