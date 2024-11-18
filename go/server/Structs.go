package handler

// Tri de l'API qu'on va traiter
type ArtistStruct struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Style        string   `json:"style"`
	Members      []string `json:"members"`
	FirstAlbum   string   `json:"firstAlbum"`
	CreationDate int      `json:"creationDate"`
	Location     []string
	// Location Locations2Struct
	Date []string
}

type LocationsStruct struct {
	Id        int      `json:"id"`
	Locations []string `json:"locations"`
}

type Locations2Struct struct {
	Index []struct {
		Id        int      `json:"id"`
		Locations []string `json:"locations"`
	} `json:"index"`
}

type DatesStruct struct {
	Id    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Dates2Struct struct {
	Index []struct {
		Id    int      `json:"id"`
		Dates []string `json:"dates"`
	} `json:"index"`
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

type Filter struct {
	Member1 bool
	Member2 bool
	Member3 bool
	Member4 bool
	Member5 bool
	Member6 bool
	Member7 bool
}
