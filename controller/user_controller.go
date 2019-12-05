package controller

import (
	util "goblog/apputil"
	"goblog/domain"
	"net/http"
)

type UserController struct {
	Store domain.UserStore
}

type CreatePage struct {
	Title          string
	SuccessMessage string
	ErrorMessage   []interface{}
}

func (uc UserController) Login(w http.ResponseWriter, r *http.Request) {
	data := CreatePage{}
	session, err := util.AppSession.SessionStore.Get(r, "goblog")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data.ErrorMessage = session.Flashes()
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	util.RenderTemplate(w, "login", data)
}

func (uc UserController) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := util.AppSession.SessionStore.Get(r, "goblog")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["user"] = domain.User{}
	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/blog/", http.StatusFound)
}

func (uc UserController) Register(w http.ResponseWriter, r *http.Request) {
	data := CreatePage{}
	data.Title = "Create User"
	util.RenderTemplate(w, "register", data)
}

func (uc UserController) CheckLogin(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	user := domain.User{}
	user.Username = username
	user.Password = password

	session, err := util.AppSession.SessionStore.Get(r, "goblog")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = uc.Store.Authenticate(user)

	if err != nil {
		session.AddFlash(err.Error())
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/login/", http.StatusFound)
		return

	}

	user.Authenticated = true
	session.Values["user"] = user

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/home/", http.StatusFound)
}

func (uc UserController) Create(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	user := domain.User{}
	user.Username = username
	user.Password = password

	_, err := uc.Store.Create(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/register/", http.StatusFound)
}
