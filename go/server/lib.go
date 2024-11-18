package handler

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

// Function to gather suggestions
func GatherSuggestions(artists []ArtistStruct) []map[string]string {
	suggestionSet := make(map[string]struct{})
	var suggestions []map[string]string

	for _, a := range artists {
		fmt.Println("LOC:", a.Location)
		log.Println("Artist:", a.Name, "Locations:", a.Location)

		AddSuggestion(suggestionSet, &suggestions, a.Name, fmt.Sprintf("/artists?id=%d", a.Id))
		AddSuggestion(suggestionSet, &suggestions, strconv.Itoa(a.CreationDate), fmt.Sprintf("/artists?id=%d", a.Id))
		AddSuggestion(suggestionSet, &suggestions, a.FirstAlbum, fmt.Sprintf("/artists?id=%d", a.Id))

		for _, m := range a.Members {
			AddSuggestion(suggestionSet, &suggestions, m, fmt.Sprintf("/artists?id=%d", a.Id))
		}
		// for _, l := range a.Location {
		// 	AddSuggestion(suggestionSet, &suggestions, l, fmt.Sprintf("/artists?id=%d", a.Id))
		// }
		// Add locations
		for _, l := range a.Location {
			AddSuggestion(suggestionSet, &suggestions, l, fmt.Sprintf("/artists?id=%d", a.Id))
		}
	}
	return suggestions
}

func AddSuggestion(set map[string]struct{}, suggestions *[]map[string]string, value string, url string) {
	lowerValue := strings.ToLower(value)
	if _, exists := set[lowerValue]; !exists && value != "" {
		set[lowerValue] = struct{}{}
		*suggestions = append(*suggestions, map[string]string{
			"Value":   value,
			"Display": value,
			"URL":     url,
		})
	}
}
