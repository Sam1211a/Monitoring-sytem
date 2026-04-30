package handlers

import (
	"html/template"
	"net/http"
)

var tmp = template.Must(template.ParseFiles("templates/dashboard.html"))

func Dashboard(w http.ResponseWriter, r *http.Request) {
	tmp.Execute(w, nil)
}
