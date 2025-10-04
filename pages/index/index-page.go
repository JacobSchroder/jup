package page_index

import (
	"html/template"
	"net/http"
)

var templates = template.Must(template.New("").ParseGlob("pages/index/*.html"))

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}
