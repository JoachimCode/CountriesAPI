package handlers

import (
	"assignment_1/utility"
	"fmt"
	"net/http"
)

const LINEBREAK = "\n"

// This method is used to handle the requests to the default endpoint.
// It will return a message to the user, suggesting some possible endpoints to try.
func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	output := "This service does not provide any functionality at this endpoint." + LINEBREAK
	output += "Please try one of the following links: " + LINEBREAK
	output += utility.INFO_PATH + utility.INFO_LIMIT_SUGGESTION + LINEBREAK
	output += utility.POPULATION_PATH + utility.POPULATION_LIMIT_SUGGESTION + LINEBREAK
	output += utility.STATUS_PATH + LINEBREAK

	_, err := fmt.Fprint(w, output)
	if err != nil {
		http.Error(w, "Error writing response.", http.StatusInternalServerError)
	}
}
