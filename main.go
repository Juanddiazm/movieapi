package main

import (
	"fmt"
	"movieapi/app"
	"movieapi/app/database"
	"net/http"
	"os"
)

// Main function to start the application
func main() {
	app := app.New()
	app.DB = &database.DB{}
	err := app.DB.Open()
	check(err)

	defer app.DB.Close()

	// Register the handler function for the / route
	http.HandleFunc("/", app.Router.ServeHTTP)
	// Start the server on port 9090
	err = http.ListenAndServe(":9090", nil)
	// Check for errors
	check(err)
}

// Check function to check for errors
func check(err error) {
	if err != nil {
		// Print the error
		fmt.Print(err.Error())
		// The program terminates immediately
		os.Exit(1)
	}
}
