package main

import (
	"net/http"
)

func overviewContent(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, overviewHtml, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
