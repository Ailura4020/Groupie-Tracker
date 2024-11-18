package handler

import (
	fetch "groupie/functions"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")

	// Vérification si l'adresse de la requète existe
	if r.URL.Path != "/" {
		// Page not found dans le cas contraire
		ErrorHandler("error404", w, nil)
	} else {
		// Stockage de l'adresse de l'API
		ArtistUrl := "https://groupietrackers.herokuapp.com/api/artists"
		AdditionnalsInfosUrl := "../static/APIs/AdditionnalsInfos.json"

		// Tri avec notre structure
		var artist []ArtistStruct
		var infos []AdditionnalsInfosStruct

		var locations2 Locations2Struct
		var dates2 Dates2Struct

		err := fetch.FetchData("https://groupietrackers.herokuapp.com/api/artists", &artist)
		err_locations := fetch.FetchData("https://groupietrackers.herokuapp.com/api/locations", &locations2)
		err_dates := fetch.FetchData("https://groupietrackers.herokuapp.com/api/dates", &dates2)

		if err_locations != nil {
			log.Println("Error fetching locations:", err_locations)
			ErrorHandler("error400", w, nil)
			return
		}
		if err_dates != nil {
			log.Println("Error fetching dates:", err_dates)
			ErrorHandler("error400", w, nil)
			return
		}
		if err != nil {
			log.Println("Error fetching artists:", err)
			http.Error(w, "Error fetching data", http.StatusInternalServerError)
			return
		}

		// Create maps for quick lookup
		locationMap := make(map[int][]string)
		dateMap := make(map[int][]string)

		for _, loc := range locations2.Index {
			locationMap[loc.Id] = loc.Locations
		}
		for _, date := range dates2.Index {
			dateMap[date.Id] = date.Dates
		}

		// Fill the Location and Date arrays in the ArtistStruct
		for i := range artist {
			a := &artist[i]
			a.Location = locationMap[a.Id]
			a.Date = dateMap[a.Id]
			artist[i] = *a // Assign the updated values back to the artist slice
		}

		// Filter the artists based on the checkbox values and search query
		var results []ArtistStruct

		for _, a := range artist {
			count := 0
			if strings.Contains(strings.ToLower(a.Name), strings.ToLower(q)) {
				results = append(results, a)
			} else if strings.Contains(strings.ToLower(strconv.Itoa(a.CreationDate)), strings.ToLower(q)) {
				results = append(results, a)
			} else if strings.Contains(strings.ToLower(a.FirstAlbum), strings.ToLower(q)) {
				results = append(results, a)
			} else {
				for _, m := range a.Members {
					if strings.Contains(strings.ToLower(m), strings.ToLower(q)) && count == 0 {
						results = append(results, a)
						count++
						break
					}
				}
				for _, location := range a.Location {
					if strings.Contains(strings.ToLower(location), strings.ToLower(q)) && count == 0 {
						results = append(results, a)
						count++
						break
					}
				}
				for _, m := range a.Date {
					if strings.Contains(strings.ToLower(m), strings.ToLower(q)) && count == 0 {
						results = append(results, a)
						count++
						break
					}
				}
				count = 0
			}
		}

		sort.Slice(results, func(i, j int) bool {
			return strings.ToLower(results[i].Name) < strings.ToLower(results[j].Name)
		})

		err_artisturl := fetch.FetchData(ArtistUrl, &artist)
		err_AdditionnalsInfos := fetch.FetchDataFromFile(AdditionnalsInfosUrl, &infos)

		// old code
		// if err_AdditionnalsInfos != nil {
		// 	log.Println(err_AdditionnalsInfos)
		// 	ErrorHandler("error400", w, nil)
		// }

		if err_AdditionnalsInfos != nil || err_artisturl != nil {
			log.Println(err_AdditionnalsInfos, err_artisturl)
			ErrorHandler("error400", w, nil)
			return
		}

		if err_artisturl != nil {
			log.Println(err_artisturl)
			ErrorHandler("error400", w, nil)
		}

		suggestions := GatherSuggestions(artist)

		// old code
		// data := map[string]interface{}{
		// 	"Artist": artist,
		// 	"AI":     infos,
		// }
		// Gather suggestions

		data := map[string]interface{}{
			"Artist": artist,
			// modified, + added suggestions
			"AI":          results,
			"Suggestions": suggestions,
		}

		// Execution de la page en envoie des données de l'API
		ErrorHandler("index", w, data)
	}
}
