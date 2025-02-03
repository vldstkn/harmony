package models

type UserStatus string

const (
	Confirmed   UserStatus = "Confirmed"
	Unconfirmed UserStatus = "Unconfirmed"
	Deleted     UserStatus = "Deleted"
)

type User struct {
	Id        int64      `db:"id"`
	CreatedAt string     `db:"created_at"`
	UpdatedAt string     `db:"updated_at"`
	Email     string     `db:"email"`
	Password  string     `db:"password"`
	Name      string     `db:"name"`
	Status    UserStatus `db:"status"`
}

type Friends struct {
	CreatedAt string `db:"created_at"`
	UserId    int64  `db:"user_id"`
	FriendId  int64  `db:"friend_id"`
}
