package router

import (
	"log"
	"net/http"
)

func New() http.Handler {
	mux := http.NewServeMux()

	//mux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.Dir("./../../swagger-ui"))))
	//mux.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
	//	http.Redirect(w, r, "/swagger/swagger.html", http.StatusFound)
	//})

	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Hello111!")
	})

	return mux
}
