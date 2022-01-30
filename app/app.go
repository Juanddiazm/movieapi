package app

import (
	"movieapi/app/database"

	"github.com/gorilla/mux"
)

// App struct
type App struct {
	Router *mux.Router
	DB     database.MovieDB
}

// Initialize App
func New() *App {
	app := &App{
		Router: mux.NewRouter(),
	}

	app.InitializeRouter()

	return app
}

// Initialize Router
func (a *App) InitializeRouter() {
	a.Router.HandleFunc("/", a.IndexHandler()).Methods("GET")
	// FindMovieByTitle
	// TODO: Review how take the title from /api/movies/{title}
	a.Router.HandleFunc("/api/movies" , a.FindMovieByTitleHandler()).Methods("GET")
	// TODO: Publish the other handlers
}
