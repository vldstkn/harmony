package models

type UserStatus string

const (
	Confirmed   UserStatus = "Confirmed"
	Unconfirmed UserStatus = "Unconfirmed"
	Deleted     UserStatus = "Deleted"
)

type User struct {
	Id        string     `db:"id"`
	CreatedAt string     `db:"created_at"`
	UpdatedAt string     `db:"updated_at"`
	Email     string     `db:"email"`
	Password  string     `db:"password"`
	Name      string     `db:"name"`
	Status    UserStatus `db:"status"`
}
