package main

import (
	"html/template"
	"net/http"
	"github.com/zmb3/spotify"
	"fmt"
	"log"
	"io/ioutil"
	"strings"
	"path/filepath"
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
	client *spotify.Client

	templateDirs = []string{"html", "data"}
	templates *template.Template
)

func getTemplates() (templates *template.Template, err error) {
	var allFiles []string
	for _, dir := range templateDirs {
		files2, _ := ioutil.ReadDir(dir)
		for _, file := range files2 {
			filename := file.Name()
			if strings.HasSuffix(filename, ".html") {
				filePath := filepath.Join(dir, filename)
				allFiles = append(allFiles, filePath)
			}
		}
	}

	templates, err = template.New("").ParseFiles(allFiles...)
	return
}

func init() {
	templates, _ = getTemplates()
}

func rootHandler(wr http.ResponseWriter, req *http.Request) {
	title := "index"

	data := map[string]interface{}{
		"title":  title,
		"header": "My Header",
		"footer": "My Footer",
	}

	err := templates.ExecuteTemplate(wr, "indexHTML", data)

	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
	}
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
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/setup", InitialSetup)
	http.HandleFunc("/callback", CompleteAuth)
	http.HandleFunc("/player/", SpotiPiControl)

	http.ListenAndServe(":8080", nil)
}