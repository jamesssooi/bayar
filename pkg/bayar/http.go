package bayar

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

// ListenAndServe creates a new Bayar server at the specified `addr` and `port`.
func ListenAndServe(addr string, port int) error {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", handleIndex)
	router.HandleFunc("/startAuthorization", handleStartGoogleAuthorization).Methods("GET")
	router.HandleFunc("/endAuthorization", handleEndGoogleAuthorization).Methods("GET")
	router.HandleFunc("/newExpense", handleExpenseCreate).Methods("POST")

	log.Printf("Starting Bayar server (v%s) on %s:%d", version, addr, port)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", addr, port), router)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "bayar server v%s", version)
}

func handleExpenseCreate(w http.ResponseWriter, r *http.Request) {
	config, _ := LoadConfig()
	var e Expense

	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	row, insertErr := e.insertIntoSpreadsheet(config.SpreadsheetID, config.SheetName)
	if insertErr != nil {
		fmt.Fprintf(w, "Error: %s", insertErr)
		return
	}

	response := struct {
		Status string `json:"status"`
		Result int    `json:"result"`
	}{"ok", row}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	if e.Label == "" {
		log.Printf("EXPENSE - Cost: %.2f  Category: %s  Row: %d", e.Cost, e.Category, row)
	} else {
		log.Printf("EXPENSE - Cost: %.2f  Category: %s  Label: %s  Row: %d", e.Cost, e.Category, e.Label, row)
	}
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
