package handler

import (
	"html/template"
	"net/http"
)

// Fonction qui s'occupe de Parse, Execute et gère les potentielles erreurs
func ErrorHandler(templ string, w http.ResponseWriter, data map[string]interface{}) {
	//Parsefile qui va s'adapter suivant le nom de la page à analyser et stocker
	page, err := template.ParseFiles("../template/" + templ + ".html")
	if err != nil {
		//error 500 si le Parsefile échoue
		error500, err3 := template.ParseFiles("../template/error500.html")
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
		error500, err3 := template.ParseFiles("../template/error500.html")
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
