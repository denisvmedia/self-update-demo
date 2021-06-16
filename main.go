package main

import (
	"log"
	"net/http"
)

var (
	Version      = "1.0.0" // semver version format
	UpdateServer = "http://localhost:8081"
)

func main() {
	must(checkRunning("", "8080"))

	go startFileServer(":8081")

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/check", handleCheck)
	http.HandleFunc("/install", handleInstall)

	log.Printf("App Version: %s", Version)
	log.Print("Now navigate to http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
