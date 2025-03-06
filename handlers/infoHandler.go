package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const countryInfoApi = "http://129.241.150.113:8080/v3.1/alpha"
const countryFlagApi = "http://129.241.150.113:3500/api/v0.1/countries/flag/images"
const citiesApi = "http://129.241.150.113:3500/api/v0.1/countries/cities"

func InfoHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleInfoGetRequest(w, r)
	default:
		http.Error(w, "Rest Method: "+r.Method+" is not supported."+"Only"+http.MethodGet+" is supported.", http.StatusMethodNotAllowed)
		return
	}
}

func handleInfoGetRequest(w http.ResponseWriter, r *http.Request) {
	countryCode := r.PathValue("countryCode")
	if len(countryCode) != 2 {
		http.Error(w, "Invalid country code. Country code should only consist of 2 characters.", http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")

	limit := 10

	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil {
			limit = parsedLimit
		}
	}

	infoResponse, err := getCountryInfo(countryCode, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		if err.Error() == "Error getting country info: 404 Not Found" {
			suggestion := "Try using a valid country code. Example: /info/{countryCode}"
			http.Error(w, suggestion, http.StatusNotFound)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)

	err = encoder.Encode(infoResponse)
	if err != nil {
		http.Error(w, "Error encoding response.", http.StatusInternalServerError)
	}
}

func getCountryInfo(countryCode string, limit int) (Info, error) {
	info, err := fetchCoreCountryInfo(countryCode)
	if err != nil {
		return Info{}, err
	}

	info.Flag, _ = getCountryFlag(countryCode)
	info.Cities, _ = fetchCities(info.Name, limit)

	if limit < len(info.Cities) {
		info.Cities = info.Cities[:limit]
	}
	return info, nil
}

func fetchCoreCountryInfo(countryCode string) (Info, error) {
	url := countryInfoApi + "/" + countryCode
	client := http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return Info{}, fmt.Errorf("Error getting country info: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Info{}, fmt.Errorf("Error getting country info: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Info{}, fmt.Errorf("Error reading API response: %v", err)
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
		return Info{}, fmt.Errorf("Error parsing API response: %v", err)
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

	info := Info{
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
		return "", fmt.Errorf("Error marshalling flag request: %v", err)
	}

	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("Error creating flag request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error getting flag: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Error getting flag: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading flag response: %v", err)
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
		return "", fmt.Errorf("Error parsing flag response: %v", err)
	}

	if apiResponse.Error {
		return "", fmt.Errorf("Error getting flag: %v", apiResponse.Msg)
	}

	return apiResponse.Data.Flag, nil
}

func fetchCities(countryFullName string, limit int) ([]string, error) {
	url := citiesApi

	payload := map[string]string{
		"country": countryFullName,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling city request: %v", err)
	}

	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("Error creating city request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error getting cities: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error getting cities: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading cities response: %v", err)
	}

	var apiResponse struct {
		Error bool     `json:"error"`
		Msg   string   `json:"msg"`
		Data  []string `json:"data"`
	}

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing cities response: %v", err)
	}

	if apiResponse.Error {
		return nil, fmt.Errorf("Error getting cities: %v", apiResponse.Msg)
	}

	return apiResponse.Data, nil
}
