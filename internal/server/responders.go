package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"text/template"
)

type errResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
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
		Code:  code,
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
	w.Write([]byte("\n"))
}

func RespondWithHTML(w http.ResponseWriter, code int, html string) {
	w.Header().Set("content-Type", "text/html")
	w.WriteHeader(code)

	_, err := w.Write([]byte(html))
	if err != nil {
		log.Printf("Failed to write HTML response: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Failed to write HTML response")
	}
}

func RespondWithTemplate(w http.ResponseWriter, code int, templatePath string, data interface{}) {
	// Parse the specified template file

	currpath, err := os.Getwd()
	if err != nil {
		log.Print("Failed to retrieve current dir path...")
		RespondWithError(w, http.StatusInternalServerError, "Error getting current directory")
		return
	}

	tmpl, err := template.ParseFiles(currpath + templatePath)
	if err != nil {
		log.Printf("Failed to parse template: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Error parsing template")
		return
	}

	// Set the content type to HTML
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(code)

	// Render the template with the provided data
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Failed to execute template: %v", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
