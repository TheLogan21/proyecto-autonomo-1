package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"ebook-system/internal/middleware"
)

type PageData struct {
	UserID int
	Error  string
	Data   interface{}
}

func render(w http.ResponseWriter, r *http.Request, page string, data interface{}) {
	t, err := template.ParseFiles("ui/html/base.tmpl", fmt.Sprintf("ui/html/pages/%s.tmpl", page))
	if err != nil {
		http.Error(w, "Error parsing templates", http.StatusInternalServerError)
		return
	}
	pd := PageData{
		UserID: middleware.GetUserID(r.Context()),
		Error:  r.URL.Query().Get("error"),
		Data:   data,
	}
	err = t.ExecuteTemplate(w, "base", pd)
	if err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
	}
}
