package api

import (
	"log"
	"net/http"
	"orderbook/gql"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
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
	srv := handler.NewDefaultServer(gql.NewExecutableSchema(gql.Config{Resolvers: &gql.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", "8080")
	log.Fatal(http.ListenAndServe(":"+"8080", nil))
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
