package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

const LINEBREAK = "\n"

func diagHandler(w http.ResponseWriter, r *http.Request) {

	output := "Requests: " + LINEBREAK
	output += "URL Path: " + r.URL.Path + LINEBREAK
	output += "Method: " + r.Method + LINEBREAK

	output += LINEBREAK + "Headers:" + LINEBREAK
	for k, v := range r.Header {
		for _, h := range v {
			output += k + ": " + h + LINEBREAK
		}
	}

	content, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body: ", err.Error())
		http.Error(w, "Error reading request body.", http.StatusInternalServerError)
		return
	}

	output += LINEBREAK + "Body:" + LINEBREAK
	output += string(content)

	_, err = fmt.Fprint(w, output)
	if err != nil {
		log.Println("Error writing response: ", err)
		http.Error(w, "Error writing response.", http.StatusInternalServerError)
	}
}
