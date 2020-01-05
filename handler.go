package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/chiefy/go-slack-utils/pkg/blockui"
	"github.com/chiefy/go-slack-utils/pkg/payload"
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

// MovieLookupHandler looks up specific movie info on TMDB and creates blocks
func MovieLookupHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var action *payload.BlockActionsPayload
	err := json.Unmarshal([]byte(r.Form.Get("payload")), &action)
	if err != nil {
		log.Printf("error decoding JSON from payload - %s", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	id, _ := strconv.Atoi(action.Actions[0].GetValue())
	movie, err := tmdbAPI.GetMovieInfo(id, map[string]string{})
	if err != nil {
		log.Printf("error looking up movie ID - %s", err)
		http.Error(w, "Bad Movie Lookup", http.StatusInternalServerError)
		return
	}
	sm := payload.NewMessagePayload("in_channel")
	ib := blockui.NewBlockImage(
		fmt.Sprintf("%s/%s/%s", tmdbImageURLBase, tmdbImageLarge, movie.PosterPath),
		movie.Title,
	)
	tb := blockui.NewBlockSection()
	tb.SetText(
		"mrkdwn",
		makeMovieMarkdown(movie),
	)
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
	res, err := tmdbAPI.SearchMovie(movieStr, map[string]string{})
	if err != nil {
		log.Println(err)
		http.Error(w, "movie lookup error", http.StatusInternalServerError)
		return
	}

	sm := payload.NewMessagePayload("ephemeral")
	sel := blockui.NewBlockSelect()

	for i, m := range res.Results {
		y := strings.Split(m.ReleaseDate, "-")
		log.Println(m.ReleaseDate)
		opt := &blockui.BlockOption{
			Text: &blockui.BlockTitleText{
				Type:  "plain_text",
				Text:  fmt.Sprintf("%s (%s)", m.Title, y[0]),
				Emoji: false,
			},
			Value: strconv.Itoa(m.ID),
		}
		sel.AddOption(opt)
		if i >= numResults-1 {
			break
		}
	}

	mb := blockui.NewBlockSectionWithSelect(sel)
	mb.SetText(
		"mrkdwn",
		fmt.Sprintf("Found %d results for *\"%s\"*:", len(res.Results), movieStr),
	)
	sm.AddBlock(mb)

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

func reverseSlice(a []string) []string {
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
	return a
}

func formatMoney(m uint32) string {
	d := strconv.Itoa(int(m))
	dSplit := reverseSlice(strings.Split(d, ""))
	f := []string{}
	for i, v := range dSplit {
		if i%3 == 0 && i != 0 {
			f = append(f, ",")
		}
		f = append(f, v)
	}
	return strings.Join(reverseSlice(f), "")
}

func makeMovieMarkdown(movie *tmdb.Movie) string {
	l := fmt.Sprintf("<https://www.imdb.com/title/%s|IMDB link>", movie.ImdbID)
	return fmt.Sprintf("*%s*\nRelease Date: %s\nBudget: $%s\n>%s\n\n%s",
		movie.Title, movie.ReleaseDate, formatMoney(movie.Budget), movie.Overview, l)
}
