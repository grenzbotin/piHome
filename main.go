package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/zmb3/spotify"
)

const (
	redirectURI        = "http://piName.local:8080/callback"
	spotify_client_id  = ""
	spotify_secret_key = ""
)

var (
	// Input
	auth  = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadCurrentlyPlaying, spotify.ScopeUserReadPlaybackState, spotify.ScopeUserModifyPlaybackState)
	ch    = make(chan *spotify.Client)
	state = "piHome"

	overviewHtml = "index.html"
	playerHtml   = "player.html"
	templates    = template.Must(template.ParseFiles("html/"+overviewHtml, "html/"+playerHtml))

	client *spotify.Client
)

type OverviewTemplateParameters struct {
	Notification string
}

type PlayerTemplateParameters struct {
	TrackName   string
	TrackArtist string
	TrackImage  string
}

func main() {
	// FileServer for template resources
	fs := http.FileServer(http.Dir("./html/res"))
	http.Handle("/res/", http.StripPrefix("/res", fs))

	// Spotify setup
	auth.SetAuthInfo(spotify_client_id, spotify_secret_key)

	// http handler
	http.HandleFunc("/", overviewContent)
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/player/", SpotiPiControl)

	go func() {
		//TODO: Authentication via UI
		url := auth.AuthURL(state)
		fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

		// wait for auth to complete
		client = <-ch

		// use the client to make calls that require authorization
		user, err := client.CurrentUser()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("You are logged in as:", user.ID)
	}()

	http.ListenAndServe(":8080", nil)
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client := auth.NewClient(tok)

	var p = OverviewTemplateParameters{
		Notification: "Login completed",
	}

	err = templates.ExecuteTemplate(w, overviewHtml, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ch <- &client
}
