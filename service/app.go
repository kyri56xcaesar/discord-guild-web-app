package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/members", membersHandler)

	certFile := os.Getenv("CERTFILE_PATH")
	keyFile := os.Getenv("KEYFILE_PATH")
	http_port := os.Getenv("HTTP_PORT")
	https_port := os.Getenv("HTTPS_PORT")
	ip := os.Getenv("IP")

	fmt.Printf("ip: %v, certFile_path: %v\nkeyFile_path: %v\nhttp_port: %v\nhttps_port: %v\n", ip, certFile, keyFile, http_port, https_port)

	if ip == "" || (https_port == "" && http_port == "") || certFile == "" || keyFile == "" {
		log.Fatalf("Required environment variables are missing")
	}

	go func() {
		fmt.Printf("Starting HTTP server at http://%v:%s\n", ip, http_port)
		err := http.ListenAndServe(":"+http_port, nil)
		if err != nil {
			log.Fatalf("HTTP server failed: %v:", err)
		}
	}()

	fmt.Printf("Starting sercure server at https://%v:%v\n", ip, https_port)
	err = http.ListenAndServeTLS(":"+https_port, certFile, keyFile, nil)
	if err != nil {
		log.Fatalf("ListenAndServeTLS failed: %v", err)
	}

}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, HTTPS world!")
}

func membersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Returing, Ranked members list")
}
