package controllers

import (
	"fmt"
	"net/http"
)

func HandleIncomingMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// To parse form data sent by Twillio
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// To extract values
	message := r.FormValue("Body")
	from := r.FormValue("From")

	fmt.Printf("Message from %s: %s\n", from, message)

	// To respond with empty 200 OK for now
	w.WriteHeader(http.StatusOK)
}
