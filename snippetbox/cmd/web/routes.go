package main

import "net/http"

// The routes() method returns a servemux containing our application routes.
func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer)) //register the file server as the handler for all URL paths that start with "/static/
	mux.HandleFunc("GET /{$}", app.home)                                // Display the home page
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)           // Display a specific snippet
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)            // Display a form for creating a new snippet
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)       // Create the new route, which is restricted to POST requests only

	return mux
}
