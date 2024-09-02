package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/sotnikea/Go_Learn/tree/main/snippetbox/internal/models"
	"github.com/sotnikea/Go_Learn/tree/main/snippetbox/ui"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates
type templateData struct {
	CurrentYear     int
	Snippet         models.Snippet
	Snippets        []models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

// Returns a nicely formatted string representation of a time.Time object
func humanDate(t time.Time) string {
	// Return the empty string if time has the zero value.
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 15:04")
}

// template.FuncMap object which acts as a lookup between the names of
// custom template functions and the functions themselves
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {

	// Initialize a new map to act as the cache
	cache := map[string]*template.Template{}

	// Use fs.Glob() to get a slice of all 'page' templates for the application
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	// Loop through the page filepaths one-by-one
	for _, page := range pages {
		// Extract the file name (like 'home.tmpl') from the full filepath
		// and assign it to the name variable
		name := filepath.Base(page)

		// Create a slice containing the filepath patterns for the templates we
		// want to parse
		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}

		// Parse the template files from the ui.Files embedded filesystem
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		// Add the template set to the map
		cache[name] = ts
	}

	// Return the map.
	return cache, nil
}
