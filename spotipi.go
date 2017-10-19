package main

import (
	"net/http"
	"log"
	"strings"
	"fmt"
)

func SpotiPiControl(w http.ResponseWriter, r *http.Request) {
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

	data := map[string]interface{}{
		"title":  "player",
		"header": "My Header - player",
		"footer": "My Footer - player",
		"trackArtist": trackArtist,
		"trackImage": trackImage,
		"trackName": trackName,
	}

	err = templates.ExecuteTemplate(w, "playerHTML", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
