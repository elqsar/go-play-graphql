package main

import (
	"fmt"
	"go-play-graphql/entity"
	"go-play-graphql/processor"
	"net/http"
	"os"

	"github.com/pressly/chi"
	"github.com/pressly/chi/render"
)

type HttpHandler struct {
	db   *entity.DB
	proc *processor.RedirectProcessor
}

func (h *HttpHandler) Redirects(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	fmt.Println("Query:", query)
	result := h.proc.ExecuteQuery(query)
	render.JSON(w, r, result)
}

func main() {
	router := chi.NewRouter()

	db := entity.NewDB()
	if err := db.HealthCheck(); err != nil {
		fmt.Println("Redis server is not available")
		os.Exit(1)
	}
	p := &processor.RedirectProcessor{Storage: db}
	handler := &HttpHandler{
		db:   db,
		proc: p,
	}
	router.Get("/redirects", handler.Redirects)

	fmt.Println("Get redirects: curl -g 'http://localhost:3000/redirects?query={redirects(offset:0,limit:30){from,to}}'")
	fmt.Printf("Get redirect: curl -g 'http://localhost:3000/redirects?query={redirect(from:%s){from,to}}'\n", "/oldUrl")
	fmt.Printf("Delete: curl -g 'http://localhost:3000/redirects?query=mutation+_{deleteRedirect(from:%s){from}}'\n", "/oldUrl")
	fmt.Printf("Create: curl -g 'http://localhost:3000/redirects?query=mutation+_{createRedirect(from:%s,to:%s){from}}'\n", "/oldUrl", "/newUrl")
	http.ListenAndServe(":3000", router)
}
