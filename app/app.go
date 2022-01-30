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
	a.Router.HandleFunc("/api/movies/{title}", a.FindMovieByTitleHandler()).Methods("GET")
	// UpdateMovieHandler
	a.Router.HandleFunc("/api/movies/{id}", a.UpdateMovieHandler()).Methods("PUT")
	// FindMoviesHandler
	a.Router.HandleFunc("/api/movies", a.FindMoviesHandler()).Methods("GET")
	// TODO: Publish the other handlers
}
