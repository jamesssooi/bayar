package bayar

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// ListenAndServe creates a new Bayar server at the specified `addr` and `port`.
func ListenAndServe(addr string, port int) {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", handleExpenseCreate).Methods("GET")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", addr, port), router))
}

func handleExpenseCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", "Hello world!")
}
