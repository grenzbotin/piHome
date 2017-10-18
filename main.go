package main

import (
	"net/http"
	"html/template"
	"fmt"
	"log"
	"strings"

	"github.com/zmb3/spotify"
)

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const redirectURI = "http://piName.local:8080/callback"

var (
	auth  = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadCurrentlyPlaying, spotify.ScopeUserReadPlaybackState, spotify.ScopeUserModifyPlaybackState)
	ch    = make(chan *spotify.Client)
	state = "piHome"

	overviewHtml = "index.html"
	playerHtml = "player.html"
	templates    = template.Must(template.ParseFiles("html/"+overviewHtml, "html/"+playerHtml))
)

type OverviewTemplateParameters struct {
	Notification     string
	TrackName		 string
	TrackArtist		 string
	TrackImage	     string
}

type Message struct {
	item string
}

func main() {
	fs := http.FileServer(http.Dir("./html/res"))
	http.Handle("/res/", http.StripPrefix("/res", fs))

	// if you didn't store your ID and secret key in the specified environment variables,
	// you can set them manually here
	auth.SetAuthInfo("clientID", "secretKey")

	// We'll want these variables sooner rather than later
	var client *spotify.Client
	var playerState *spotify.PlayerState

	/*http.Handle("/", http.FileServer(http.Dir("/home/pi/piServer/html")))
	http.ListenAndServe(":3000", nil)*/

	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var p = OverviewTemplateParameters{
			Notification:             "Hello",
		}
		err := templates.ExecuteTemplate(w, overviewHtml, p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/player/", func(w http.ResponseWriter, r *http.Request) {
		action := strings.TrimPrefix(r.URL.Path, "/player/")
		fmt.Println("Got request for:", action)
		var err error
		switch action {
		case "play":
			err = client.Play()
		case "pause":
			err = client.Pause()
		case "next":
			err = client.Next()
		case "previous":
			err = client.Previous()
		case "shuffle":
			playerState.ShuffleState = !playerState.ShuffleState
			err = client.Shuffle(playerState.ShuffleState)
		}
		if err != nil {
			log.Print(err)
		}

		currently, err := client.PlayerCurrentlyPlaying()

		trackName := currently.Item.Name
		trackArtist := currently.Item.Artists[0].Name
		trackImage := currently.Item.Album.Images[1].URL


		if err != nil {
			log.Print(err)
		}
		/*
				w.Header().Set("Content-Type", "text/html")
				fmt.Fprint(w, html)
	*/
		var p = OverviewTemplateParameters{
			TrackName:		  trackName,
			TrackArtist:	  trackArtist,
			TrackImage:       trackImage,
		}

		err = templates.ExecuteTemplate(w, playerHtml, p)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	})


	go func() {
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

		/*
		playerState, err = client.PlayerState()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Found your %s (%s)\n", playerState.Device.Type, playerState.Device.Name)
		*/
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
	/*w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "Login Completed!"+html)*/

	var p = OverviewTemplateParameters{
		Notification:             "Login completed",
	}
	err = templates.ExecuteTemplate(w, overviewHtml, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ch <- &client
}

