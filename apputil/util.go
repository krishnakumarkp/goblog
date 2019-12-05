package apputil

import (
	"fmt"
	"goblog/domain"
	"net/http"
	"text/template"

	"github.com/gorilla/sessions"
)

func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles("templates/" + tmpl + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetUser(s *sessions.Session) domain.User {
	val := s.Values["user"]
	fmt.Println(val)
	var user = domain.User{}
	user, ok := val.(domain.User)
	if !ok {
		return domain.User{Authenticated: false}
	}
	return user
}
