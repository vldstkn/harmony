package notifications

import "harmony/pkg/db"

type Repository struct {
	DB *db.DB
}

func NewRepository(database *db.DB) *Repository {
	return &Repository{
		DB: database,
	}
}

func (repo *Repository) Save() (string, error) {
	return "", nil
}
func (repo *Repository) GetById() (string, error) {
	return "", nil
}
func (repo *Repository) GetMany(userId int64) (string, error) {
	return "", nil
}
