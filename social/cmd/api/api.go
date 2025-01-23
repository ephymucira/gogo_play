package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/chi/v5"
)



type application struct{
	config config
}

type config struct{
	addr string
}
func (app *application) mount() http.Handler {
	// mux := http.NewServeMux()

	// mux.HandleFunc("GET /v1/health",app.healthCheckHandler)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	//recoverer which helps recover from a paniv attack
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	// r.Get("/", func(w http.ResponseWriter, r *http.Request){
	// 	w.Write([]byte("welcome"))
	// })

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
	})

	return r
}

func (app *application) run(mux http.Handler) error {

	

	srv := &http.Server{
		Addr: app.config.addr,
		Handler: mux,
		WriteTimeout: time.Second *30,
		ReadTimeout: time.Second *10,
		IdleTimeout: time.Minute,

	}
	log.Printf("Server has started at %s", app.config.addr)

	return srv.ListenAndServe()
}