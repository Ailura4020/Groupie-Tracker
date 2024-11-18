package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Tri de l'API qu'on va traiter
type ArtistStruct struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Style        string   `json:"style"`
	Members      []string `json:"members"`
	FirstAlbum   string   `json:"firstAlbum"`
	CreationDate int      `json:"creationDate"`
}

type LocationsStruct struct {
	Id        int      `json:"id"`
	Locations []string `json:"locations"`
}

type DatesStruct struct {
	Id    int      `json:"id"`
	Dates []string `json:"dates"`
}

type AdditionnalsInfosStruct struct {
	Id        int    `json:"id"`
	Image     string `json:"image"`
	Name      string `json:"name"`
	Style     string `json:"style"`
	Biography string `json:"bio"`
}

type RelationsStruct struct {
	Locations string
	Dates     string
}

func main() {
	server := &http.Server{
		Addr:              ":8081",          //adresse du server (le port choisi est à titre d'exemple) // listes des handlers
		ReadHeaderTimeout: 10 * time.Second, // temps autorisé pour lire les headers
		WriteTimeout:      10 * time.Second, // temps maximum d'écriture de la réponse
		IdleTimeout:       60 * time.Second, // temps maximum entre deux rêquetes
		MaxHeaderBytes:    1 << 20,          // 1 MB // maxinmum de bytes que le serveur va lire
	}
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/artists", artistHandler)
	http.HandleFunc("/search", searchHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	//Vérification si l'adresse de la requète existe
	if r.URL.Path != "/" {
		//Page not found dans le cas contraire
		errorHandler("error404", w, nil)
	} else {
		//Stockage de l'adresse de l'API
		ArtistUrl := "https://groupietrackers.herokuapp.com/api/artists"
		AdditionnalsInfosUrl := "static/APIs/AdditionnalsInfos.json"

		//Tri avec notre structure
		var artist []ArtistStruct
		var infos []AdditionnalsInfosStruct

		err_artisturl := FetchData(ArtistUrl, &artist)
		err_AdditionnalsInfos := FetchDataFromFile(AdditionnalsInfosUrl, &infos)

		if err_AdditionnalsInfos != nil {
			log.Println(err_AdditionnalsInfos)
			errorHandler("error400", w, nil)
		}
		if err_artisturl != nil {
			log.Println(err_artisturl)
			errorHandler("error400", w, nil)
		}

		data := map[string]interface{}{
			"Artist": artist,
			"AI":     infos,
		}
		//Execution de la page en envoie des données de l'API
		errorHandler("index", w, data)
	}
}

func artistHandler(w http.ResponseWriter, r *http.Request) {
	//Stockage de l'id de l'artiste sur lequel on va cliquer
	id := r.URL.Query().Get("id")

	//On récupère la partie de l'API qui nous interesse pour relation et artist
	artisturl := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/artists/%s", id)
	locationsurl := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/locations/%s", id)
	datesurl := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/dates/%s", id)
	infosurl := "static/APIs/AdditionnalsInfos.json"

	// Stockage des structures dans des variables
	var artist ArtistStruct
	var locations LocationsStruct
	var dates DatesStruct
	var infos []AdditionnalsInfosStruct

	// Fetch des datas qu'on va recuperer
	err_artist := FetchData(artisturl, &artist)
	err_locations := FetchData(locationsurl, &locations)
	err_dates := FetchData(datesurl, &dates)
	err_infos := FetchDataFromFile(infosurl, &infos)
	infos_id, err_Atoi := strconv.Atoi(id)

	//Vérification que les differents FetchData ont bien fonctionnés
	if err_artist != nil {
		log.Println(err_artist)
		errorHandler("error400", w, nil)
	}
	if err_locations != nil {
		log.Println(err_locations)
		errorHandler("error400", w, nil)
	}
	if err_dates != nil {
		log.Println(err_dates)
		errorHandler("error400", w, nil)
	}
	if err_infos != nil {
		log.Println(err_infos)
		errorHandler("error400", w, nil)
	}
	if err_Atoi != nil {
		log.Println(err_Atoi)
		return
	}

	// Creation d'une structure data qui regroupe les 3 structures
	var relations []RelationsStruct
	for i := 0; i < len(locations.Locations); i++ {
		relation := RelationsStruct{
			Locations: Capitalize(strings.ReplaceAll(strings.ReplaceAll(locations.Locations[i], "-", " - "), "_", " ")),
			Dates:     strings.ReplaceAll(strings.ReplaceAll(dates.Dates[i], "*", ""), "-", "/"),
		}
		relations = append(relations, relation)
	}

	data := map[string]interface{}{
		"Artist":    artist,
		"Relations": relations,
		"AI":        infos[infos_id-1],
	}
	errorHandler("artist", w, data)

}

// Fonction qui s'occupe de Parse, Execute et gère les potentielles erreurs
func errorHandler(templ string, w http.ResponseWriter, data map[string]interface{}) {
	//Parsefile qui va s'adapter suivant le nom de la page à analyser et stocker
	page, err := template.ParseFiles("template/" + templ + ".html")
	if err != nil {
		//error 500 si le Parsefile échoue
		error500, err3 := template.ParseFiles("template/error500.html")
		if err3 != nil {
			http.Error(w, "Error 500", http.StatusInternalServerError)
			return
		}
		err4 := error500.Execute(w, data)
		if err4 != nil {
			http.Error(w, "Error 500", http.StatusInternalServerError)
			return
		}
		return
	}
	err2 := page.Execute(w, data)
	if err2 != nil {
		error500, err3 := template.ParseFiles("template/error500.html")
		if err3 != nil {
			http.Error(w, "Error 500", http.StatusInternalServerError)
			return
		}
		err4 := error500.Execute(w, data)
		if err4 != nil {
			http.Error(w, "Error 500", http.StatusInternalServerError)
			return
		}
		return
	}
}

func FetchData(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("HTTP request error: %v", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code: %d", resp.StatusCode)
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	err = json.NewDecoder(resp.Body).Decode(target)
	if err != nil {
		log.Printf("JSON decode error: %v", err)
		return err
	}
	return nil
}

func FetchDataFromFile(filePath string, target interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return err
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(target)
	if err != nil {
		log.Printf("JSON decode error: %v", err)
		return err
	}
	return nil
}

func Capitalize(s string) string {
	var result string
	IsNewWord := true
	for _, l := range s {
		alph := (l >= 'a' && l <= 'z') || (l >= 'A' && l <= 'Z') || (l >= '0' && l <= '9')
		if alph {
			if IsNewWord {
				if l >= 'a' && l <= 'z' {
					l = l + -32
				}
				IsNewWord = false
			} else {
				if l >= 'A' && l <= 'Z' {
					l = l + 32
				}
			}
		} else {
			IsNewWord = true
		}
		result += string(l)
	}
	return result
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		http.Error(w, "Search query is empty", http.StatusBadRequest)
		return
	}

	// Fetch data from API or database based on the search query
	var artist []ArtistStruct
	var location []RelationsStruct

	err := FetchData("https://groupietrackers.herokuapp.com/api/artists", &artist)
	err1 := FetchData("https://groupietrackers.herokuapp.com/api/relation", &location)
	fmt.Println(location)
	if err != nil {
		http.Error(w, "Error fetching data", http.StatusInternalServerError)
		return
	}
	if err1 != nil {
		http.Error(w, "Error fetching data", http.StatusInternalServerError)
		return
	}
	// Filter the artists based on the search query
	var results []ArtistStruct
	for _, a := range artist {
		if strings.Contains(strings.ToLower(a.Name), strings.ToLower(q)) {
			results = append(results, a)
		} else if strings.Contains(strings.ToLower(strconv.Itoa(a.CreationDate)), strings.ToLower(q)) {
			results = append(results, a)
		} else if strings.Contains(strings.ToLower(a.FirstAlbum), strings.ToLower(q)) {
			results = append(results, a)
		} else {
			for _, m := range a.Members {
				if strings.Contains(strings.ToLower(m), strings.ToLower(q)) {
					results = append(results, a)
					break
				}
			}
		}
	}
	// Filter the relation based on the search query
	var test []RelationsStruct
	for _, b := range location {
		if strings.Contains(strings.ToLower(b.Locations), strings.ToLower(q)) {
			test = append(test, b)
			break
		}
	}

	// Return the search results as JSON
	// json.NewEncoder(w).Encode(results)
	data := map[string]interface{}{
		"Artist":   results,
		"Location": test,
	}
	errorHandler("search", w, data)
}
