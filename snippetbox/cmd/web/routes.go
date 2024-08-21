package main

import (
	"net/http"

	"github.com/justinas/alice" // New import
)

// The routes() method returns a servemux containing our application routes.
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer)) //register the file server as the handler for all URL paths that start with "/static/

	mux.HandleFunc("GET /{$}", app.home)                          // Display the home page
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)     // Display a specific snippet
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)      // Display a form for creating a new snippet
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost) // Create the new route, which is restricted to POST requests only

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	// Return the 'standard' middleware chain followed by the servemux.
	return standard.Then(mux)
}
