package domain

import "database/sql"

type Blog struct {
	ID           int
	Title        string
	Details      string
	Description  string
	CreatedDate  string
	CreatedTime  string
	ModifiedDate string
	ModifiedTime string
	Public       string
	Image        sql.NullString
}

type User struct {
	ID            int
	Username      string
	Password      string
	Authenticated bool
}

type BlogStore interface {
	Create(Blog) (Blog, error)
	Update(string, Blog) error
	Delete(string) error
	GetById(string) (Blog, error)
	GetAll() ([]Blog, error)
}

type UserStore interface {
	Create(User) (User, error)
	Authenticate(User) (User, error)
}
