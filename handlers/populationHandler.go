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

const populationApi = "http://129.241.150.113:3500/api/v0.1/countries/population"
const isoApi = "http://129.241.150.113:3500/api/v0.1/countries/iso"

// this method is used to handle the requests to the /population endpoint.
func PopulationHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handlePopulationGetRequest(w, r)
	default:
		utility.StatusWriter(w, http.StatusMethodNotAllowed, "Rest Method: "+r.Method+" is not supported. Only "+http.MethodGet+" is supported.")
		return
	}
}

// this method is used to handle the GET requests to the /population endpoint.
func handlePopulationGetRequest(w http.ResponseWriter, r *http.Request) {
	countryCode := r.PathValue("countryCode")
	if len(countryCode) != 2 {
		utility.StatusWriter(w, http.StatusBadRequest, "Invalid country code. Country code should only consist of 2 characters.")
		return
	}

	limitStr := r.URL.Query().Get("limit")

	iso3, err := getIso3(countryCode)
	if iso3 == "Bad request" {
		utility.StatusWriter(w, http.StatusBadRequest, "Error: did not find country code")
		return
	}

	if err != nil {
		utility.StatusWriter(w, http.StatusInternalServerError, err.Error())
		return
	}

	var startYear, endYear int
	startYear = 0
	endYear = 0

	if limitStr != "" {
		years := strings.Split(limitStr, "-")
		if len(years) != 2 {
			utility.StatusWriter(w, http.StatusBadRequest, "Invalid year range. Should be in the format: /population/{countryCode}?limit=2010-2020")
			return
		}
		if len(years) == 2 {
			startYear, err = strconv.Atoi(years[0])
			if err != nil {
				utility.StatusWriter(w, http.StatusBadRequest, "Invalid start year. Should be in the format: /population/{countryCode}?limit=2010-2020")
				return
			}
			endYear, err = strconv.Atoi(years[1])
			if err != nil {
				utility.StatusWriter(w, http.StatusBadRequest, "Invalid end year. Should be in the format: /population/{countryCode}?limit=2010-2020")
				return
			}
		}
	}

	populationResponse, err := getPopulationInfo(iso3, startYear, endYear)
	if err != nil {
		if strings.Contains(err.Error(), "Invalid year") {
			utility.StatusWriter(w, http.StatusBadRequest, err.Error())
			return
		} else {
			utility.StatusWriter(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)

	err = encoder.Encode(populationResponse)
	if err != nil {
		utility.StatusWriter(w, http.StatusInternalServerError, "Error encoding response")
		return
	}
}

// this method is used to get the iso3 code for a given country code.
func getIso3(countryCode string) (string, error) {
	url := isoApi

	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("Error getting iso data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Error getting iso data: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading API response: %v", err)
	}

	var apiResponse struct {
		Error bool   `json:"error"`
		Msg   string `json:"msg"`
		Data  []struct {
			Name string `json:"name"`
			Iso2 string `json:"iso2"`
			Iso3 string `json:"iso3"`
		} `json:"data"`
	}

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return "", fmt.Errorf("Error unmarshalling API response: %v", err)
	}

	countryCode = strings.ToUpper(countryCode)

	for _, country := range apiResponse.Data {
		if country.Iso2 == countryCode {
			return strings.ToLower(country.Iso3), nil
		}
	}

	return "Bad request", nil
}

// this method is used to get the population information for a given country code.
func getPopulationInfo(countryCode string, startYear, endYear int) ([]utility.YearlyPopulation, error) {
	url := populationApi

	payload := map[string]interface{}{
		"iso3": countryCode,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return []utility.YearlyPopulation{}, fmt.Errorf("Error marshalling population request: %v", err)
	}

	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return []utility.YearlyPopulation{}, fmt.Errorf("Error creating population request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return []utility.YearlyPopulation{}, fmt.Errorf("Error getting population: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []utility.YearlyPopulation{}, fmt.Errorf("Error getting population: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []utility.YearlyPopulation{}, fmt.Errorf("Error reading population response: %v", err)
	}

	var apiResponse struct {
		Error bool   `json:"error"`
		Msg   string `json:"msg"`
		Data  struct {
			Country          string `json:"country"`
			Code             string `json:"code"`
			Iso3             string `json:"iso3"`
			PopulationCounts []struct {
				Year       int `json:"year"`
				Population int `json:"value"`
			} `json:"populationCounts"`
		} `json:"data"`
	}

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return []utility.YearlyPopulation{}, fmt.Errorf("Error parsing population response: %v", err)
	}

	if apiResponse.Error {
		return []utility.YearlyPopulation{}, fmt.Errorf("Error getting population: %v", apiResponse.Msg)
	}

	var populationData []utility.YearlyPopulation

	if startYear == 0 && endYear == 0 {
		for _, year := range apiResponse.Data.PopulationCounts {
			populationData = append(populationData, utility.YearlyPopulation{Year: year.Year, Population: year.Population})
		}
		return populationData, nil
	}

	if startYear > endYear {
		return []utility.YearlyPopulation{}, fmt.Errorf("Invalid year range. Start year should be less than or equal to end year.")
	}

	for _, year := range apiResponse.Data.PopulationCounts {
		if year.Year >= startYear && year.Year <= endYear {
			populationData = append(populationData, utility.YearlyPopulation{Year: year.Year, Population: year.Population})
		}
	}
	return populationData, nil
}
