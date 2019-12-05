package mysqlstore

import (
	"errors"
	"goblog/domain"
)

type BlogStore struct {
	Store DataStore
}

func (bs BlogStore) GetById(id string) (domain.Blog, error) {
	var blog domain.Blog

	sqlStatement := "select id, details, title, description, image, date_posted, time_posted, date_edited, time_edited, public FROM blog where id = ?"
	rows, err := bs.Store.Db.Query(sqlStatement, id)
	if err != nil {
		return blog, err
	}
	blognotfound := true
	for rows.Next() {
		blognotfound = false
		err = rows.Scan(&blog.ID, &blog.Details, &blog.Title, &blog.Description, &blog.Image, &blog.CreatedDate, &blog.CreatedTime, &blog.ModifiedDate, &blog.ModifiedTime, &blog.Public)
		if err != nil {
			return blog, err
		}
	}
	if blognotfound {
		return blog, errors.New("User not found")
	}
	return blog, nil
}
func (bs BlogStore) GetAll() ([]domain.Blog, error) {
	var blogs []domain.Blog
	var blog domain.Blog

	sqlStatement := `SELECT id, details, title, description, image, date_posted, time_posted, date_edited, time_edited, public FROM blog ORDER BY id DESC;`
	rows, err := bs.Store.Db.Query(sqlStatement)
	if err != nil {
		return blogs, errors.New("could not retrive blogs")
	}

	for rows.Next() {

		blog = domain.Blog{}
		err := rows.Scan(&blog.ID, &blog.Details, &blog.Title, &blog.Description, &blog.Image, &blog.CreatedDate, &blog.CreatedTime, &blog.ModifiedDate, &blog.ModifiedTime, &blog.Public)
		if err != nil {
			return nil, err
		}
		blogs = append(blogs, blog)
	}
	return blogs, nil
}

func (bs BlogStore) Create(blog domain.Blog) (domain.Blog, error) {
	sqlStatement := "INSERT INTO blog(title, description, details, date_posted, time_posted, date_edited, time_edited, public) VALUES (?,?,?,?,?,?,?,?)"
	insForm, err := bs.Store.Db.Prepare(sqlStatement)
	if err != nil {
		return domain.Blog{}, err
	}
	_, ierr := insForm.Exec(blog.Title, blog.Description, blog.Details, blog.CreatedDate, blog.CreatedTime, blog.ModifiedDate, blog.ModifiedTime, blog.Public)
	if ierr != nil {
		return domain.Blog{}, ierr
	}
	return blog, nil
}

func (bs BlogStore) Delete(id string) error {
	if id == "" {
		return errors.New("Id is not passed")
	}
	sqlStatement := "DELETE FROM blog WHERE id = ?"
	_, err := bs.Store.Db.Exec(sqlStatement, id)
	if err != nil {
		return err
	}
	return nil
}

func (bs BlogStore) Update(id string, blog domain.Blog) error {

	updateQuery := "UPDATE blog SET title=?, description=?, details=?, public=?, date_edited=?, time_edited=? WHERE id=?"

	_, err := bs.Store.Db.Exec(updateQuery, blog.Title, blog.Description, blog.Details, blog.Public, blog.ModifiedDate, blog.ModifiedTime, id)
	if err != nil {
		return err
	}
	return nil
}
