package main

import (
	"fmt"
	"net/http"

	"webcrawlergo/crawler"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/leiloes", crawler.GetLeiloesHandler).Methods("GET")

	fmt.Println("Ta on e roteando. Use a rota: /leiloes?url={URL}")
	http.ListenAndServe(":8000", router)
}
