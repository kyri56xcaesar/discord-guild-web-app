package server

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
)

type errResponse struct {
	Error string `json:"error"`
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {

	if code > 499 {
		log.Println("Responding with 5XX error: ", msg)
	}

	if code > 399 {
		log.Println("Responding with 4XX error: ", msg)
	}

	RespondWithJSON(w, code, errResponse{
		Error: msg,
	})

}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON payload: %+v", payload)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func IsAlphanumeric(s string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return re.MatchString(s)
}
