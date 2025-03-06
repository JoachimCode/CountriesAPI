package handlers

type Info struct {
	Name       string   `json:"name"`
	Continents []string `json:"continents"`
	Population int      `json:"population"`
	Languages  []string `json:"languages"`
	Borders    []string `json:"borders"`
	Flag       string   `json:"flag"`
	Capital    string   `json:"capital"`
	Cities     []string `json:"cities"`
}

type YearlyPopulation struct {
	Year       int `json:"year"`
	Population int `json:"population"`
}

type Population struct {
	Values []YearlyPopulation `json:"values"`
}

type IsoInformation struct {
	Name string `json:"name"`
	Iso2 string `json:"Iso2"`
	Iso3 string `json:"Iso3"`
}

type StatusInformation struct {
	CountriesNowApi  string
	RestCountriesApi string
	Version          string
	Uptime           string
}
