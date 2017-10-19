package main

import "net/http"
import "log"

var p = SetupParameters{
	SetupSpotifyNotification: "",
	SetupSpotifyLink:"",
}

func InitialSetup(w http.ResponseWriter, r *http.Request) {
	createSpotiPiLink()

	err := templates.ExecuteTemplate(w, setupHtml, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func createSpotiPiLink() {
	url := auth.AuthURL(state)
	p.SetupSpotifyNotification = "Please log in to Spotify by clicking the following page in your browser:"
	p.SetupSpotifyLink = url
}

func CompleteAuth(w http.ResponseWriter, r *http.Request) {
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

	var p = SetupParameters{
		SetupSpotifyNotification: "Login completed",
	}

	err = templates.ExecuteTemplate(w, overviewHtml, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ch <- &client
}