package controller

import (
	util "goblog/apputil"
	"goblog/domain"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type BlogController struct {
	Store domain.BlogStore
}

//var sessionStore *sessions.CookieStore

func (bc BlogController) ListAll(w http.ResponseWriter, r *http.Request) {

	blogs, err := bc.Store.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		//http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	util.RenderTemplate(w, "index", blogs)
}
func (bc BlogController) Home(w http.ResponseWriter, r *http.Request) {

	type data struct {
		User           domain.User
		Blogs          []domain.Blog
		SuccessMessage []interface{}
	}

	var err error

	formData := data{}
	session, err := util.AppSession.SessionStore.Get(r, "goblog")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := util.GetUser(session)

	//loggedin := session.Values["loggedin"]
	//user := domain.User{}
	//user, ok := val.(User)

	if auth := user.Authenticated; !auth {
		//if loggedin == "yes" {
		session.AddFlash("You don't have access!")
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/login/", http.StatusFound)
		return
	}

	formData.Blogs, err = bc.Store.GetAll()
	formData.User = user

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		//http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}

	formData.SuccessMessage = session.Flashes()
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.RenderTemplate(w, "home", formData)
}

func (bc BlogController) Add(w http.ResponseWriter, r *http.Request) {

	blog := domain.Blog{}

	blog.Title = r.FormValue("title")
	blog.Description = r.FormValue("description")
	blog.Details = r.FormValue("details")
	public := r.FormValue("public")
	if public == "yes" {
		blog.Public = "yes"
	} else {
		blog.Public = "no"
	}

	todaysDate := time.Now().Format("02-01-2006")
	t := time.Now()

	time := t.Format("15:04")

	blog.CreatedDate = todaysDate
	blog.CreatedTime = time
	blog.ModifiedDate = todaysDate
	blog.ModifiedTime = time

	//fmt.Printf("%v", public)

	//fmt.Printf("%+v\n", blog)

	_, err := bc.Store.Create(blog)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/home/", http.StatusFound)
}

func (bc BlogController) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := bc.Store.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/home/", http.StatusFound)
}

func (bc BlogController) Edit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var err error
	type data struct {
		Blog  domain.Blog
		Blogs []domain.Blog
	}

	formData := data{}
	formData.Blog, err = bc.Store.GetById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	formData.Blogs, err = bc.Store.GetAll()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	util.RenderTemplate(w, "edit", formData)
}

func (bc BlogController) Update(w http.ResponseWriter, r *http.Request) {

	blog := domain.Blog{}
	id := r.FormValue("id")
	blog.Title = r.FormValue("title")
	blog.Description = r.FormValue("description")
	blog.Details = r.FormValue("details")
	public := r.FormValue("public")
	if public == "yes" {
		blog.Public = "yes"
	} else {
		blog.Public = "no"
	}

	todaysDate := time.Now().Format("02-01-2006")
	t := time.Now()
	time := t.Format("15:04")

	blog.ModifiedDate = todaysDate
	blog.ModifiedTime = time

	err := bc.Store.Update(id, blog)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	session, err := util.AppSession.SessionStore.Get(r, "goblog")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.AddFlash("Blog Updated")
	err = session.Save(r, w)

	http.Redirect(w, r, "/home/", http.StatusFound)
}
