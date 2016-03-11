package main

import (
	"log"
	"os"
)

func main() {
	dns_service := os.Getenv("GODYN_SERVICE")
	dns_user := os.Getenv("GODYN_USER")
	dns_pw := os.Getenv("GODYN_PW")
	if dns_service == "" || dns_user == "" || dns_pw == "" {
		log.Printf("Please ensure GODYN_SERVICE, GODYN_USER, and GODYN_PW are set")
		return
	}

	log.Printf("Using %v service", dns_service)
}
