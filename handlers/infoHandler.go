package handlers

import (
	"assignment_1/utility"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const countryInfoApi = "http://129.241.150.113:8080/v3.1/alpha"
const countryFlagApi = "http://129.241.150.113:3500/api/v0.1/countries/flag/images"
const citiesApi = "http://129.241.150.113:3500/api/v0.1/countries/cities"

// This method is used to handle the requests to the /info endpoint.
// It will return information about a country, including its flag and a list of cities.
func InfoHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleInfoGetRequest(w, r)
	default:
		utility.StatusWriter(w, http.StatusMethodNotAllowed, "Rest Method: "+r.Method+" is not supported. Only "+http.MethodGet+" is supported.")
		return
	}
}

// This method is used to handle the GET requests to the /info endpoint.
func handleInfoGetRequest(w http.ResponseWriter, r *http.Request) {
	countryCode := r.PathValue("countryCode")
	if len(countryCode) != 2 {
		utility.StatusWriter(w, http.StatusBadRequest, "Invalid country code. Country code should only consist of 2 characters.")
		return
	}

	limitStr := r.URL.Query().Get("limit")

	limit := 10

	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil {
			limit = parsedLimit
		} else {
			utility.StatusWriter(w, http.StatusBadRequest, "Invalid limit. Should be in the format: /info/{countryCode}?limit=10. Will show 10 cities by default.")
		}
	}

	infoResponse, err := getCountryInfo(countryCode, limit)

	if err != nil {
		if strings.Contains(err.Error(), "Bad request") {
			utility.StatusWriter(w, http.StatusBadRequest, "Error: did not find country code")
			return
		} else {
			utility.StatusWriter(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)

	err = encoder.Encode(infoResponse)
	if err != nil {
		utility.StatusWriter(w, http.StatusInternalServerError, "Error encoding response")
		return
	}
}

// This method is used to fetch information about a country, including its flag and a list of cities.
func getCountryInfo(countryCode string, limit int) (utility.Info, error) {
	info, err := fetchCoreCountryInfo(countryCode)
	if err != nil {
		return utility.Info{}, err
	}

	info.Flag, err = getCountryFlag(countryCode)
	if err != nil {
		return utility.Info{}, err
	}

	info.Cities, err = fetchCities(info.Name, limit)
	if err != nil {
		return utility.Info{}, err
	}

	if limit < len(info.Cities) {
		info.Cities = info.Cities[:limit]
	}
	return info, nil
}

// This method is used to fetch core information about a country.
func fetchCoreCountryInfo(countryCode string) (utility.Info, error) {
	url := countryInfoApi + "/" + countryCode
	client := http.Client{Timeout: 5 * time.Second}
	startTime.Local().Location()
	resp, err := client.Get(url)
	if err != nil {
		return utility.Info{}, fmt.Errorf("error getting country info: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusBadRequest {
		return utility.Info{}, fmt.Errorf("Bad request")
	} else if resp.StatusCode != http.StatusOK {
		return utility.Info{}, fmt.Errorf("error getting country info: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return utility.Info{}, fmt.Errorf("error reading API response: %v", err)
	}

	var apiResponse []struct {
		Name struct {
			Common string `json:"common"`
		} `json:"name"`
		Continents []string          `json:"continents"`
		Population int               `json:"population"`
		Languages  map[string]string `json:"languages"`
		Borders    []string          `json:"borders"`
		Capital    []string          `json:"capital"`
	}

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return utility.Info{}, fmt.Errorf("error parsing API response: %v", err)
	}

	countryData := apiResponse[0]

	var languageList []string
	for _, lang := range countryData.Languages {
		languageList = append(languageList, lang)
	}

	var capital string
	if len(countryData.Capital) > 0 {
		capital = countryData.Capital[0]
	} else {
		capital = "Unknown"
	}

	info := utility.Info{
		Name:       countryData.Name.Common,
		Continents: countryData.Continents,
		Population: countryData.Population,
		Languages:  languageList,
		Borders:    countryData.Borders,
		Capital:    capital,
		Cities:     []string{},
		Flag:       "",
	}

	return info, nil
}

// This method is used to fetch the flag of a country.
func getCountryFlag(countryCode string) (string, error) {
	url := countryFlagApi

	if len(countryCode) > 2 {
		countryCode = countryCode[:2]
	}

	payload := map[string]string{
		"iso2": countryCode,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error marshalling flag request: %v", err)
	}

	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating flag request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error getting flag: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusBadRequest {
		return "", fmt.Errorf("Bad request")
	} else if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error getting country info: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading flag response: %v", err)
	}

	var apiResponse struct {
		Error bool   `json:"error"`
		Msg   string `json:"msg"`
		Data  struct {
			Flag string `json:"flag"`
		} `json:"data"`
	}

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return "", fmt.Errorf("error parsing flag response: %v", err)
	}

	if apiResponse.Error {
		return "", fmt.Errorf("error getting flag: %v", apiResponse.Msg)
	}

	return apiResponse.Data.Flag, nil
}

// This method is used to fetch a list of cities in a country.
func fetchCities(countryFullName string, limit int) ([]string, error) {
	url := citiesApi

	payload := map[string]string{
		"country": countryFullName,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshalling city request: %v", err)
	}

	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating city request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error getting cities: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusBadRequest {
			return nil, fmt.Errorf("Bad request")
		} else if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("error getting cities: %v", resp.Status)
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading cities response: %v", err)
	}

	var apiResponse struct {
		Error bool     `json:"error"`
		Msg   string   `json:"msg"`
		Data  []string `json:"data"`
	}

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing cities response: %v", err)
	}

	if apiResponse.Error {
		return nil, fmt.Errorf("error getting cities: %v", apiResponse.Msg)
	}

	return apiResponse.Data, nil
}
