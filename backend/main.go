package main

import (
	"log"
	"net/http"
)

func main() {
	printBanner()

	// Load configuration
	log.Println("Using configuration: ...")

	// Serve static content
	httpRegisterAssets()

	// Start http server
	http.ListenAndServe(":7340", nil)
}
