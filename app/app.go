package app

import "github.com/gorilla/mux"

// App struct
type App struct {
	Router *mux.Router
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
	a.Router.HandleFunc("/", nil).Methods("GET")
}
