package app

import (
	"fmt"
	"net/http"
)

// IndexHandler is the handler for the index page
func (a *App) IndexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	}
}