package main

import (
	"log"
	"net/http"
	"path/filepath"
)

func startFileServer(addr string) {
	dir, err := filepath.Abs("./updates")
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(addr, http.FileServer(http.Dir(dir))))
}