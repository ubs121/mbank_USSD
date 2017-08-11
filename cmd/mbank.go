package main

import (
	ussd "../stk"
	"log"
	"net/http"
)

func main() {
	// service binding
	mux := http.NewServeMux()
	ussd.RegisterService(mux)

	log.Println("hub is started...")
	//http.ListenAndServeTLS(":4002", "cert.pem", "key.pem", mux)
	http.ListenAndServe(":8080", mux)

	log.Println("hub is stopped.")
}
