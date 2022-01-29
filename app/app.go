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
}
