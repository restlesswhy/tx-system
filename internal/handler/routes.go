package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func LoadRoutes(r *mux.Router, handler *Handler) *mux.Router {
	api := r.PathPrefix("/api").Subrouter()
	api.Use(Logging)
	api.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("PONG!")) })

	user := api.PathPrefix("/user").Subrouter()
	user.HandleFunc("/create", handler.CreateUser).Methods(http.MethodPost)

	userBalance := user.PathPrefix("/balance").Subrouter()
	userBalance.HandleFunc("/update", handler.UpdateBalance).Methods(http.MethodPut)

	return api
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s\n", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
