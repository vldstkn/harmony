package account

import (
	"errors"
	"harmony/internal/models"
	"harmony/pkg/db"
)

type Repository struct {
	DB *db.DB
}

func NewRepository(db *db.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func (repo *Repository) GetById(id string) *models.User {
	var user models.User
	err := repo.DB.Get(&user, `SELECT * FROM users WHERE id=$1`, id)
	if err != nil {
		return nil
	}
	return &user
}
func (repo *Repository) GetByEmail(email string) *models.User {
	var user models.User
	err := repo.DB.Get(&user, `SELECT * FROM users WHERE email=$1`, email)
	if err != nil {
		return nil
	}
	return &user
}

func (repo *Repository) UserNotExists(email, name string) error {
	var user models.User
	err := repo.DB.Get(&user, `SELECT name, email FROM users WHERE email=$1 OR name=$2`, email, name)
	if err != nil {
		return nil
	}
	if user.Email == email {
		return errors.New("email")
	} else {
		return errors.New("name")
	}
}

func (repo *Repository) Create(user models.User) (string, error) {
	var id string
	err := repo.DB.QueryRow(`INSERT INTO users (email, password, name) values ($1,$2,$3) RETURNING id`,
		user.Email, user.Password, user.Name).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}
func (repo *Repository) DeleteUnique(field, value string) error {
	_, err := repo.DB.Exec(`DELETE FROM users WHERE $1=$2`, field, value)
	return err
}
func (repo *Repository) SearchByName(name string) []models.User {
	var users []models.User
	err := repo.DB.Select(&users, `SELECT * FROM users WHERE name LIKE '$1%'`, name)
	if err != nil {
		return nil
	}
	return users
}
func (repo *Repository) UpdateById(field, value, id string) error {
	_, err := repo.DB.Exec(`UPDATE users SET status=$1 WHERE id=$2`, value, id)
	return err
}
