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
