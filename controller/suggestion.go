package controller

import (
	"encoding/json"
	"net/http"
)

type SuggestionController struct{}

func NewSuggestionController() *SuggestionController {
	return &SuggestionController{}
}

func (s *SuggestionController) SuggestedSubredditsHandler(w http.ResponseWriter, r *http.Request) {
	var request []string
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: Implement collaborative recommending system
	suggestions := []string{
		"t5_2qh0u", // pics
		"t5_2sbq3", // EarthPorn
		"t5_2scjs", // CityPorn
		"t5_2r1tc", // itookapicture
	}

	writeJSON(w, suggestions, http.StatusOK)
}
