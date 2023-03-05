package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/oschvr/bookings/pkg/config"
	"github.com/oschvr/bookings/pkg/handlers"
	"github.com/oschvr/bookings/pkg/render"
	"log"
	"net/http"
	"time"
)

const host = "localhost"
const port = "8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	// Is production
	app.IsProduction = false

	// Session init
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.Secure = app.IsProduction

	app.Session = session

	// Template cache init
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatalln("Error creating template cache")
	}
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: routes(&app),
	}

	fmt.Println(fmt.Sprintf("Starting application on port %s", port))
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalln("error starting server: ", err)
	}
}
