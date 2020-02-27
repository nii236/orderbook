package api

import (
	"net/http"

	"github.com/go-chi/chi"
)

type Server struct {
}

func New() {
	r := chi.NewRouter()
	r.Post("/orderbooks/{orderbook_id}/add", add)
	r.Post("/orderbooks/{orderbook_id}/remove", remove)
	r.Get("/orderbooks/{orderbook_id}", orderbookDetails)
	r.Get("/orderbooks", orderbookList)

}

func add(w http.ResponseWriter, r *http.Request) {
	return
}
func remove(w http.ResponseWriter, r *http.Request) {
	return
}
func orderbookDetails(w http.ResponseWriter, r *http.Request) {
	return
}
func orderbookList(w http.ResponseWriter, r *http.Request) {
	return
}
