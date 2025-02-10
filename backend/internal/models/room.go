package models

type RoomRole int

const (
	Owner     RoomRole = 2
	Moderator RoomRole = 1
	Member    RoomRole = 0
	Banned    RoomRole = -1
)

type Room struct {
	Id           int64  `db:"id"`
	Name         string `db:"name"`
	CreatorId    int64  `db:"creator_id"`
	CreatedAt    string `db:"created_at"`
	UpdatedAt    string `db:"updated_at"`
	Participants []int64
}

type RoomMember struct {
	RoomId   int64    `db:"room_id"`
	UserId   int64    `db:"user_id"`
	JoinedAt string   `db:"joined_at"`
	Role     RoomRole `db:"role"`
}

func (r RoomRole) String() string {
	switch r {
	case Banned:
		return "BANNED"
	case Member:
		return "MEMBER"
	case Moderator:
		return "MODERATOR"
	case Owner:
		return "OWNER"
	default:
		return "UNKNOWN"
	}
}
