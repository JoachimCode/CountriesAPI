package handlers

import (
	"fmt"
	"net/http"
)

func DefaultHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")
	output := "This service does not provide any functionality at this endpoint." + LINEBREAK
	output += "Please try: " + LINEBREAK
	output += INFO_PATH + LINEBREAK

	_, err := fmt.Fprint(w, output)
	if err != nil {
		http.Error(w, "Error writing response.", http.StatusInternalServerError)
	}
}
