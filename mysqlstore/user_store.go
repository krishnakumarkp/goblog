package mysqlstore

import (
	"errors"
	"goblog/domain"
)

type UserStore struct {
	Store DataStore
}

func (userStore UserStore) Create(user domain.User) (domain.User, error) {
	rows, err := userStore.Store.Db.Query("select count(*) as usercount from users where username = ?", user.Username)
	if err != nil {
		return domain.User{}, err
	}
	var usercount int64
	for rows.Next() {
		rows.Scan(&usercount)
	}
	if usercount > 0 {
		return domain.User{}, errors.New("this username is taken")
	}
	sqlStatement := "INSERT INTO users(username, password) VALUES (?,?)"
	insForm, err := userStore.Store.Db.Prepare(sqlStatement)
	if err != nil {
		return domain.User{}, err
	}
	_, ierr := insForm.Exec(user.Username, user.Password)
	if ierr != nil {
		return domain.User{}, ierr
	}
	return user, nil
}

func (userStore UserStore) Authenticate(user domain.User) (domain.User, error) {

	var username string
	var password string
	sqlStatement := "SELECT username,password from users where username = ?"
	rows, err := userStore.Store.Db.Query(sqlStatement, user.Username)
	if err != nil {
		return domain.User{}, err
	}
	usernotfound := true
	for rows.Next() {
		usernotfound = false
		err = rows.Scan(&username, &password)
		if err != nil {
			return domain.User{}, err
		}
	}
	if usernotfound {
		return domain.User{}, errors.New("User not found")
	}
	if user.Username == username && user.Password == password {
		return user, nil
	}
	return domain.User{}, errors.New("Invalid username password")

}
