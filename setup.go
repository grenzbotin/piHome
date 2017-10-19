package main

import "net/http"
import "log"

func InitialSetup(w http.ResponseWriter, r *http.Request) {
	url := auth.AuthURL(state)

	data := map[string]interface{}{
		"title":  "Start",
		"SetupSpotifyNotification" : "Please log in to Spotify by clicking the following page in your browser:",
		"SetupSpotifyLink": url,
		"header": "My Header - Setup",
		"footer": "My Footer - Setup",
	}

	err := templates.ExecuteTemplate(w, "setupHTML", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

	data := map[string]interface{}{
		"title":  "Start",
		"header": "My Header - Index",
		"footer": "My Footer - Index",
	}

	err = templates.ExecuteTemplate(w, "indexHTML", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ch <- &client
}