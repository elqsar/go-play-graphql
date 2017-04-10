package main

import (
	"github.com/pressly/chi"
	"net/http"
	"fmt"
	"github.com/pressly/chi/render"
	"go-play-graphql/processor"
	"go-play-graphql/entity"
)


func main() {
	router := chi.NewRouter()
	router.Get("/redirects", Redirects)

	fmt.Println("Get redirects: curl -g 'http://localhost:3000/redirects?query={redirects(offset:0,limit:30){from,to}}'")
	fmt.Printf("Get redirect: curl -g 'http://localhost:3000/redirects?query={redirect(from:%s){from,to}}'\n", "/oldUrl")
	fmt.Printf("Delete: curl -g 'http://localhost:3000/redirects?query=mutation+_{deleteRedirect(from:%s){from}}'\n", "/oldUrl")
	fmt.Printf("Create: curl -g 'http://localhost:3000/redirects?query=mutation+_{createRedirect(from:%s,to:%s){from}}'\n", "/oldUrl", "/newUrl")
	http.ListenAndServe(":3000", router)
}

func Redirects(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	fmt.Println("Query:", query)
	db := entity.NewDB()
	p := &processor.RedirectProcessor{Storage: db}
	result := p.ExecuteQuery(query)
	render.JSON(w, r, result)
}