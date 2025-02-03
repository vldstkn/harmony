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

func (repo *Repository) GetById(id int64) *models.User {
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

func (repo *Repository) Create(user models.User) (int64, error) {
	var id int64
	err := repo.DB.QueryRow(`INSERT INTO users (email, password, name) values ($1,$2,$3) RETURNING id`,
		user.Email, user.Password, user.Name).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}
func (repo *Repository) DeleteById(id int64) error {
	_, err := repo.DB.Exec(`DELETE FROM users WHERE id=$2`, id)
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
func (repo *Repository) UpdateStatusById(value string, id int64) error {
	_, err := repo.DB.Exec(`UPDATE users SET status=$1 WHERE id=$2`, value, id)
	return err
}

func (repo *Repository) FindByName(id int64, name string) []models.User {
	var users []models.User
	name += "%"
	err := repo.DB.Select(&users, `SELECT * 
																				FROM users 
																		    WHERE name LIKE $2 AND id!=$1`, id, name)
	if err != nil {
		return nil
	}
	return users
}
func (repo *Repository) AddFriend(userId, friendId int64) error {
	if userId > friendId {
		userId, friendId = friendId, userId
	}
	_, err := repo.DB.Exec(`INSERT INTO friendships (user_id1, user_id2) 
																 VALUES ($1, $2)`, userId, friendId)
	return err

}
func (repo *Repository) DeleteFriend(userId, friendId int64) error {
	_, err := repo.DB.Exec(`DELETE FROM friendships 
       													 WHERE (user_id2=$1 AND user_id1=$2) OR (user_id2=$2 AND user_id1=$1)`, userId, friendId)
	return err
}
func (repo *Repository) FindFriendsByName(userId int64, name string) []models.User {
	var users []models.User
	name += "%"
	err := repo.DB.Select(&users, `SELECT u.name, u.email, u.created_at, u.id
																				FROM users u
																				JOIN friendships f ON u.id = f.user_id1 OR u.id = f.user_id2
																				WHERE (f.user_id1 = $1 OR f.user_id2 = $1) AND u.id != $1 AND u.name LIKE $2;`, userId, name)
	if err != nil {
		return nil
	}
	return users
}
