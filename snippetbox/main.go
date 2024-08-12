package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// Define a home handler function which writes a byte slice containing
// "Hello from Snippetbox" as the response body.
func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	w.Write([]byte("Hello from Snippetbox"))
}

// Add a snippetView handler function.
func snippetView(w http.ResponseWriter, r *http.Request) {
	// Extract the value of the id wildcard from the request using r.PathValue()
	// and try to convert it to an integer using the strconv.Atoi() function. If
	// it can't be converted to an integer, or the value is less than 1, we
	// return a 404 page not found response.
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

// Add a snippetCreate handler function.
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

// Add a snippetCreatePost handler function.
func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Save a new snippet..."))
}

// Test snippet
func testMarishka(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Marishka lapunya, naykrawa kycunya"))
}

// Test snippet2
func testMamulya(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Мамуля, я тебя люблю!"))
}

func main() {
	// Use the http.NewServeMux() function to initialize a new servemux
	mux := http.NewServeMux()

	// Register handler functions
	mux.HandleFunc("GET /{$}", home)                          // Display the home page
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)     // Display a specific snippet
	mux.HandleFunc("GET /snippet/create", snippetCreate)      // Display a form for creating a new snippet
	mux.HandleFunc("POST /snippet/create", snippetCreatePost) // Create the new route, which is restricted to POST requests only.

	mux.HandleFunc("GET /marishka", testMarishka)
	mux.HandleFunc("GET /mama", testMamulya)

	// Print a log message to say that the server is starting.
	log.Print("starting server on :4000")

	// Use the http.ListenAndServe() function to start a new web server. We pass in
	// two parameters: the TCP network address to listen on (in this case ":4000")
	// and the servemux we just created. If http.ListenAndServe() returns an error
	// we use the log.Fatal() function to log the error message and exit. Note
	// that any error returned by http.ListenAndServe() is always non-nil.
	err := http.ListenAndServe("0.0.0.0:4000", mux)
	log.Fatal(err)
}
