package handler

import (
	"fmt"
	fetch "groupie/functions"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	//Stockage de l'id de l'artiste sur lequel on va cliquer
	id := r.URL.Query().Get("id")

	//On récupère la partie de l'API qui nous interesse pour relation et artist
	artisturl := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/artists/%s", id)
	locationsurl := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/locations/%s", id)
	datesurl := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/dates/%s", id)
	infosurl := "../static/APIs/AdditionnalsInfos.json"

	// Stockage des structures dans des variables
	var artist ArtistStruct
	var locations LocationsStruct
	var dates DatesStruct
	var infos []AdditionnalsInfosStruct

	// Fetch des datas qu'on va recuperer
	err_artist := fetch.FetchData(artisturl, &artist)
	err_locations := fetch.FetchData(locationsurl, &locations)
	err_dates := fetch.FetchData(datesurl, &dates)
	err_infos := fetch.FetchDataFromFile(infosurl, &infos)
	infos_id, err_Atoi := strconv.Atoi(id)

	//Vérification que les differents FetchData ont bien fonctionnés
	if err_artist != nil {
		log.Println(err_artist)
		ErrorHandler("error400", w, nil)
	}
	if err_locations != nil {
		log.Println(err_locations)
		ErrorHandler("error400", w, nil)
	}
	if err_dates != nil {
		log.Println(err_dates)
		ErrorHandler("error400", w, nil)
	}
	if err_infos != nil {
		log.Println(err_infos)
		ErrorHandler("error400", w, nil)
	}
	if err_Atoi != nil {
		log.Println(err_Atoi)
		return
	}

	// Creation d'une structure data qui regroupe les 3 structures
	var relations []RelationsStruct
	for i := 0; i < len(locations.Locations); i++ {
		relation := RelationsStruct{
			Locations: fetch.Capitalize(strings.ReplaceAll(strings.ReplaceAll(locations.Locations[i], "-", " - "), "_", " ")),
			Dates:     strings.ReplaceAll(strings.ReplaceAll(dates.Dates[i], "*", ""), "-", "/"),
		}
		relations = append(relations, relation)
	}

	data := map[string]interface{}{
		"Artist":    artist,
		"Relations": relations,
		"AI":        infos[infos_id-1],
	}
	ErrorHandler("artist", w, data)
}
