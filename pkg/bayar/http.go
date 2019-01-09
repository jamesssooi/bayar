package bayar

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

// ListenAndServe creates a new Bayar server at the specified `addr` and `port`.
func ListenAndServe(addr string, port int) {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/startAuthorization", handleStartGoogleAuthorization).Methods("GET")
	router.HandleFunc("/endAuthorization", handleEndGoogleAuthorization).Methods("GET")
	router.HandleFunc("/", handleExpenseCreate).Methods("GET")

	log.Printf("Starting Bayar server on %s:%d", addr, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", addr, port), router))
}

func handleExpenseCreate(w http.ResponseWriter, r *http.Request) {
	config, _ := LoadConfig()
	fmt.Fprintf(w, "Hello world! %s", config.GoogleConfigurationFilename)
}

func handleStartGoogleAuthorization(w http.ResponseWriter, r *http.Request) {
	googlecfg, err := loadGoogleClientConfig()
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	url := googlecfg.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Fprintf(w, "<html>")
	fmt.Fprintf(w, "Get your authorization code <a href='%s'>here</a>", url)
	fmt.Fprintf(w, "<br><hr><br>")
	fmt.Fprintf(w, "<form method='GET' action='/endAuthorization'>")
	fmt.Fprintf(w, "Then enter your code here: <input type='text' id='code' name='code'/><input type='submit' value='Submit'/>")
	fmt.Fprintf(w, "</form>")
	fmt.Fprintf(w, "</html>")
}

func handleEndGoogleAuthorization(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	if q.Get("code") == "" {
		fmt.Fprintf(w, "Error: No authorization code provide.")
		return
	}

	_, err := processAuthorizationCode(r.FormValue("code"))
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
	}

	fmt.Fprintf(w, "You are successfully authorized!")
}
