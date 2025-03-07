package handlers

import (
	"assignment_1/utility"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

const restApi = "http://129.241.150.113:8080/"
const countryApi = "http://129.241.150.113:8080/"

var startTime time.Time

// This method is used to handle the requests to the /status endpoint.
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleStatusGetRequest(w, r)
	default:
		utility.StatusWriter(w, http.StatusMethodNotAllowed, "Rest Method: "+r.Method+" is not supported. Only "+http.MethodGet+" is supported.")
		return
	}
}

// This method is used to handle the GET requests to the /status endpoint.
func handleStatusGetRequest(w http.ResponseWriter, r *http.Request) {
	statusRestApi, err := checkRestStatus()
	if err != nil {
		utility.StatusWriter(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")

	statusCountryApi, err := checkCountryApiStatus()
	if err != nil {
		utility.StatusWriter(w, http.StatusInternalServerError, err.Error())
		return
	}

	secondsSinceStart := time.Since(startTime).Seconds()
	uptime := "time in seconds since the last re/start of your service: " + strconv.FormatFloat(secondsSinceStart, 'f', -1, 64)

	statusInformation := utility.StatusInformation{
		CountriesNowApi:  statusCountryApi,
		RestCountriesApi: statusRestApi,
		Version:          "v1.0",
		Uptime:           uptime,
	}

	encoder := json.NewEncoder(w)

	err = encoder.Encode(statusInformation)
	if err != nil {
		utility.StatusWriter(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// This method is used to check the status of the REST API.
// It sends a HEAD request to the REST API and returns the status of the response.
func checkRestStatus() (string, error) {
	url := restApi

	client := http.Client{
		Timeout: time.Second * 5,
	}

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	return resp.Status, nil
}

// This method is used to check the status of the country API.
// It sends a HEAD request to the country API and returns the status of the response.
func checkCountryApiStatus() (string, error) {
	url := countryApi

	client := http.Client{
		Timeout: time.Second * 5,
	}

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	return resp.Status, nil
}

// This method is used to set the start time of the service.
// It is used to calculate the uptime of the service.
func SetStartTime(givenTime time.Time) {
	startTime = givenTime
}
