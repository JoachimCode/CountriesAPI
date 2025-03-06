package handlers

import (
	"fmt"
	"net/http"
)

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	output := "This service does not provide any functionality at this endpoint." + LINEBREAK
	output += "Please try one of the following links: " + LINEBREAK
	output += INFO_PATH + INFO_LIMIT_SUGGESTION + LINEBREAK
	output += POPULATION_PATH + POPULATION_LIMIT_SUGGESTION + LINEBREAK
	output += STATUS_PATH + LINEBREAK

	_, err := fmt.Fprint(w, output)
	if err != nil {
		http.Error(w, "Error writing response.", http.StatusInternalServerError)
	}
}
