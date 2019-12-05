package main

import (
	"encoding/gob"
	util "goblog/apputil"
	"goblog/controller"
	"goblog/domain"
	"goblog/mysqlstore"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/spf13/viper"
)

func main() {

	config := mysqlstore.Config{
		Host:     util.AppConfig.DBHost,
		Port:     util.AppConfig.DBPort,
		User:     util.AppConfig.DBUser,
		Password: util.AppConfig.DBPassword,
		Database: util.AppConfig.Database,
	}
	// Creates a Mysql DB instance
	dbstore, err := mysqlstore.New(config)
	if err != nil {
		panic(err)
	}

	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	sessionStore := sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	util.AppSession = util.NewSession(sessionStore)

	userStore := mysqlstore.UserStore{
		dbstore,
	}

	blogStore := mysqlstore.BlogStore{
		dbstore,
	}

	userController := controller.UserController{
		Store: userStore,
	}

	blogController := controller.BlogController{
		Store: blogStore,
	}

	r := mux.NewRouter()
	r.HandleFunc("/login/", userController.Login).Methods("GET")
	r.HandleFunc("/logout/", userController.Logout).Methods("GET")
	r.HandleFunc("/blog/", blogController.ListAll).Methods("GET")
	r.HandleFunc("/register/", userController.Register).Methods("GET")
	r.HandleFunc("/createuser/", userController.Create).Methods("POST")
	r.HandleFunc("/checklogin/", userController.CheckLogin).Methods("POST")
	r.HandleFunc("/home/", blogController.Home).Methods("GET")
	r.HandleFunc("/add/", blogController.Add).Methods("POST")
	r.HandleFunc("/delete/{id}", blogController.Delete).Methods("GET")
	r.HandleFunc("/edit/{id}", blogController.Edit).Methods("GET")
	r.HandleFunc("/update/", blogController.Update).Methods("POST")
	//r.HandleFunc("/text", wrapHandler(textHandler)).Methods("POST")
	//r.HandleFunc("/text/{hash}", wrapHandler(textHashHandler)).Methods("GET")

	//http.HandleFunc("/login/", userController.Login)
	log.Fatal(http.ListenAndServe(":8080", r))
}

func init() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	util.AppConfig.DBHost = viper.GetString("mysql.Host")
	util.AppConfig.DBPort = viper.GetString("mysql.Port")
	util.AppConfig.DBUser = viper.GetString("mysql.User")
	util.AppConfig.DBPassword = viper.GetString("mysql.Password")
	util.AppConfig.Database = viper.GetString("mysql.Database")

	gob.Register(domain.User{})

}
