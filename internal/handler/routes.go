package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func LoadRoutes(r *mux.Router, handler *Handler) *mux.Router {
	api := r.PathPrefix("/api").Subrouter()
	api.Use(Logging)
	api.HandleFunc("/update", handler.Update).Methods("POST")
	api.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("PONG!")) })

	return api
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s\n", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
