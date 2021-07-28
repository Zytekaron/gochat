package old_server

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"gochat/src/events"
	"net/http"
)

func newRouter() chi.Router {
	router := chi.NewRouter()

	cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // fixme prod danger CORS Policy
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	})

	router.Use(middleware.RequestID, middleware.RealIP, middleware.Logger, middleware.Recoverer)

	router.Get("/", serveHome)
	router.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(server, w, r)
	})

	router.Route("/messages", func(r chi.Router) {
		//r.Get("/", wrap(events.MessageGetBatch))
		//r.Get("/{id}", wrap(events.MessageGet))
		r.Post("/", func(w http.ResponseWriter, req *http.Request) {
			events.MessageSend(server, cmds, w, req)
		})
		r.Put("/{id}", wrap(events.MessageEdit))
		r.Delete("/{id}", wrap(events.MessageDelete))
	})
	return router
}