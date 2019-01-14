package bayar

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"

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

	if tokenErr := checkVerificationToken(r.Header); tokenErr != nil {
		sendErrorResponse(w, tokenErr.Error(), 400)
		return
	}

	if decodeErr := json.NewDecoder(r.Body).Decode(&e); decodeErr != nil {
		sendErrorResponse(w, decodeErr.Error(), 400)
		return
	}

	row, insertErr := e.InsertIntoSpreadsheet(config.SpreadsheetID, config.SheetName)
	if insertErr != nil {
		sendErrorResponse(w, insertErr.Error(), 500)
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
		sendErrorResponse(w, err.Error(), 500)
		return
	}

	url := googlecfg.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Fprintf(w, "<html>")
	fmt.Fprintf(w, "Get your authorization code <a href='%s' target='_blank'>here</a>", url)
	fmt.Fprintf(w, "<br><hr><br>")
	fmt.Fprintf(w, "<form method='GET' action='/endAuthorization'>")
	fmt.Fprintf(w, "Then enter your code here: <input type='text' id='code' name='code'/><input type='submit' value='Submit'/>")
	fmt.Fprintf(w, "</form>")
	fmt.Fprintf(w, "</html>")
}

func handleEndGoogleAuthorization(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	if q.Get("code") == "" {
		sendErrorResponse(w, "no authorization code provided", 400)
		return
	}

	_, err := processAuthorizationCode(r.FormValue("code"))
	if err != nil {
		sendErrorResponse(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, "You are successfully authorized!")
}

func sendErrorResponse(w http.ResponseWriter, msg string, code int) {
	d := struct {
		Ok    bool   `json:"ok"`
		Error string `json:"error"`
	}{false, msg}
	m, _ := json.MarshalIndent(d, "", "  ")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprint(w, string(m))
}

func checkVerificationToken(headers http.Header) error {
	re := regexp.MustCompile(`^Bearer (.+)$`)
	match := re.FindStringSubmatch(headers.Get("Authorization"))
	if match == nil {
		return errors.New("missing verification token")
	}

	token := match[1]
	config, _ := LoadConfig()
	if token != config.VerificationToken {
		return errors.New("invalid verification token")
	}

	return nil
}
