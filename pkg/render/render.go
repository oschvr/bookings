package render

import (
	"bytes"
	"github.com/oschvr/bookings/pkg/config"
	"github.com/oschvr/bookings/pkg/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var app *config.AppConfig

// NewTemplates sets the config for template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// AddDefaultData adds default data to all templates
func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// RenderTemplate renders html template
func RenderTemplate(w http.ResponseWriter, tmpl string, data *models.TemplateData) {

	// If dev mode, rebuild templateCache, if not, read from it
	var tc map[string]*template.Template
	if app.UseCache {
		// get template cache from app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Println("Error requesting template from cache")
	}

	buf := new(bytes.Buffer)

	td := AddDefaultData(data)

	_ = t.Execute(buf, td)

	// render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println("Error rendering template", err)
		log.Println(err)
	}
}

// CreateTemplateCache creates a template cachec
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all the files name *.page.tmpl.html from ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl.html")
	if err != nil {
		log.Println("Error going through pages", err)
		return myCache, err
	}

	// range through the found pages, ending in *.tmpl.html
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			log.Printf("Error parsing file %s", name)
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl.html")

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl.html")
			if err != nil {
				log.Println("Error parsing glob", err)
				return myCache, err
			}
		}
		if err != nil {
			log.Println("Error finding matches", err)
			return myCache, err
		}

		myCache[name] = ts
	}

	return myCache, nil
}
