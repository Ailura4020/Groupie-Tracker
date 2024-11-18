package handler

import (
	"fmt"
	fetch "groupie/functions"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func FiltersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	//Initializing my structs
	var artist []ArtistStruct
	var locations2 Locations2Struct
	var dates2 Dates2Struct

	//Fetching the datas on my Structs
	err := fetch.FetchData("https://groupietrackers.herokuapp.com/api/artists", &artist)
	err_locations := fetch.FetchData("https://groupietrackers.herokuapp.com/api/locations", &locations2)
	err_dates := fetch.FetchData("https://groupietrackers.herokuapp.com/api/dates", &dates2)

	//Checking for fetch errors
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
		ErrorHandler("error400", w, nil)
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

	//Stock the form values into a variable
	formValues := r.URL.Query()

	//Create a boolean array that will contains checkboxes values
	filter := Filter{}
	members := formValues["members"]
	for _, value := range members {
		switch value {
		case "1":
			filter.Member1 = true
		case "2":
			filter.Member2 = true
		case "3":
			filter.Member3 = true
		case "4":
			filter.Member4 = true
		case "5":
			filter.Member5 = true
		case "6":
			filter.Member6 = true
		case "7":
			filter.Member7 = true
		}
	}

	//Filtering the struct with numbers of members values from form
	var results []ArtistStruct
	for _, a := range artist {
		if len(a.Members) == 1 && filter.Member1 {
			results = append(results, a)
		}
		if len(a.Members) == 2 && filter.Member2 {
			results = append(results, a)
		}
		if len(a.Members) == 3 && filter.Member3 {
			results = append(results, a)
		}
		if len(a.Members) == 4 && filter.Member4 {
			results = append(results, a)
		}
		if len(a.Members) == 5 && filter.Member5 {
			results = append(results, a)
		}
		if len(a.Members) == 6 && filter.Member6 {
			results = append(results, a)
		}
		if len(a.Members) == 7 && filter.Member7 {
			results = append(results, a)
		}
	}

	//Get the min and max values for the creation
	mincreation := formValues.Get("mincreation")
	maxcreation := formValues.Get("maxcreation")

	//Convert those values into ints
	creationmin, errmin := strconv.Atoi(mincreation)
	creationmax, errmax := strconv.Atoi(maxcreation)

	//If one of those isn't an int then end the programm
	if errmin != nil || errmax != nil {
		fmt.Println(errmin)
		fmt.Println(errmax)
		return
	}

	//If the struct is empty then it takes all the values, aka artist
	if len(results) == 0 {
		results = artist
	}

	//Filter the struct with Creation Date form values
	var creationfilter []ArtistStruct
	for _, a := range results {
		if a.CreationDate >= creationmin && a.CreationDate <= creationmax {
			creationfilter = append(creationfilter, a)
		}
	}

	//Get the min and max values for the date
	mindate := formValues.Get("mindate")
	maxdate := formValues.Get("maxdate")

	//Convert those values into ints
	datemin, errmin := strconv.Atoi(mindate)
	datemax, errmax := strconv.Atoi(maxdate)

	//If one of those isn't an int then end the programm
	if errmin != nil || errmax != nil {
		fmt.Println(errmin)
		fmt.Println(errmax)
		return
	}

	//If the struct is empty then it takes all the values, aka artist
	if len(creationfilter) == 0 {
		creationfilter = artist
	}

	//Filter the struct with First Album form values
	var firstalbumfilter []ArtistStruct
	for _, a := range creationfilter {
		for i := datemin; i <= datemax; i++ {
			if strings.Contains(a.FirstAlbum, strconv.Itoa(i)) {
				firstalbumfilter = append(firstalbumfilter, a)
			}
		}
	}

	//If the struct is empty then it takes all the values, aka artist
	if len(firstalbumfilter) == 0 {
		firstalbumfilter = artist
	}

	//Filter the previous struct with the values from the search barr for the location
	var locationfilter []ArtistStruct
	loc := formValues.Get("location")
	for _, a := range firstalbumfilter {
		for _, location := range a.Location {
			if strings.Contains(loc, ",") {
				if strings.Contains(loc, ", ") {
					tab := strings.Split(loc, ", ")
					state := tab[0]
					country := tab[1]
					if strings.Contains(strings.ToLower(location), strings.ToLower(state)) && strings.Contains(strings.ToLower(location), strings.ToLower(country)) {
						locationfilter = append(locationfilter, a)
						break
					}
				} else {
					tab := strings.Split(loc, ",")
					state := tab[0]
					country := tab[1]
					if strings.Contains(strings.ToLower(location), strings.ToLower(state)) && strings.Contains(strings.ToLower(location), strings.ToLower(country)) {
						locationfilter = append(locationfilter, a)
						break
					}
				}
			} else {
				if strings.Contains(strings.ToLower(location), strings.ToLower(loc)) {
					locationfilter = append(locationfilter, a)
					break
				}
			}
		}
	}

	//If the struct is empty then it takes all the values, aka artist
	if len(locationfilter) == 0 {
		locationfilter = artist
	}

	//Sending the struct after filtering it
	data := map[string]interface{}{
		"AI": locationfilter,
	}
	ErrorHandler("filters", w, data)
}
