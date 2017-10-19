package main

import (
	"html/template"
	"net/http"
	"github.com/zmb3/spotify"
	"fmt"
	"log"
)

const (
	redirectURI        = "http://soda.local:8080/callback"
	spotify_client_id  = "clientId"
	spotify_secret_key = "secretKey"
)

var (
	// Input
	auth  = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadCurrentlyPlaying, spotify.ScopeUserReadPlaybackState, spotify.ScopeUserModifyPlaybackState)
	ch    = make(chan *spotify.Client)
	state = "piHome"

	overviewHtml = "index.html"
	setupHtml = "setup.html"
	playerHtml   = "player.html"
	templates    = template.Must(template.ParseFiles("html/"+overviewHtml, "html/"+playerHtml, "html/"+setupHtml))

	client *spotify.Client
)

type OverviewTemplateParameters struct {
	Welcome string
}

type SetupParameters struct {
	SetupSpotifyNotification string
	SetupSpotifyLink string
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

	go func() {
		// wait for auth to complete
		client = <-ch

		// use the client to make calls that require authorization
		user, err := client.CurrentUser()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("You are logged in as:", user.ID)
	}()

	// http handler
	http.HandleFunc("/", overviewContent)
	http.HandleFunc("/setup", InitialSetup)
	http.HandleFunc("/callback", CompleteAuth)
	http.HandleFunc("/player/", SpotiPiControl)

	http.ListenAndServe(":8080", nil)
}