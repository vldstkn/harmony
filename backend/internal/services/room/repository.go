package room

import (
	"fmt"
	"harmony/internal/models"
	"harmony/pkg/db"
	"strings"
)

type Repository struct {
	DB *db.DB
}

func NewRepository(database *db.DB) *Repository {
	return &Repository{
		DB: database,
	}
}
func (repo *Repository) Create(creatorId int64, name string) (int64, error) {
	tx, err := repo.DB.Beginx()
	if err != nil {
		return -1, err
	}
	var roomId int64
	err = tx.QueryRow(`INSERT INTO rooms (name, creator_id) VALUES ($1, $2) RETURNING id`, name, creatorId).Scan(&roomId)
	if err != nil {
		tx.Rollback()
		return -1, err
	}
	_, err = tx.Exec(`INSERT INTO rooms_members (room_id, user_id, role) VALUES ($1, $2, 2)`, roomId, creatorId)
	if err != nil {
		tx.Rollback()
		return -1, err
	}
	err = tx.Commit()
	if err != nil {
		return -1, err
	}
	return roomId, nil
}
func (repo *Repository) Delete(roomId int64) error {
	_, err := repo.DB.Exec(`DELETE FROM rooms WHERE id=$1`, roomId)
	return err
}
func (repo *Repository) AddUsers(roomId int64, usersId []int64) error {
	var values []string
	args := []interface{}{roomId}
	for i, userId := range usersId {
		values = append(values, fmt.Sprintf("($1, $%d)", i+2))
		args = append(args, userId)
	}
	query := fmt.Sprintf("INSERT INTO rooms_members (room_id, user_id) VALUES %s ON CONFLICT (room_id, user_id) DO NOTHING;", strings.Join(values, ", "))
	_, err := repo.DB.Exec(query, args...)
	return err
}
func (repo *Repository) RemoveUsers(roomId int64, userIds []int64) error {
	var placeholders []string
	args := []interface{}{roomId}
	for _, userId := range userIds {
		placeholders = append(placeholders, "$"+fmt.Sprint(len(args)+1))
		args = append(args, userId)
	}
	query := fmt.Sprintf("DELETE FROM rooms_members WHERE room_id=$1 AND user_id IN (%s);", strings.Join(placeholders, ", "))
	_, err := repo.DB.Exec(query, args...)
	return err
}
func (repo *Repository) GetRoomsByUserId(userId int64) []models.Room {
	var rooms []models.Room
	err := repo.DB.Select(&rooms, `SELECT r.id, r.creator_id, r.name, r.created_at 
																 FROM rooms_members rm
														     JOIN rooms r ON r.id=rm.room_id AND rm.user_id=$1
														     `, userId)
	if err != nil {
		return nil
	}
	return rooms
}

func (repo *Repository) GetRoomById(id int64) *models.Room {
	var room models.Room
	err := repo.DB.Get(&room, `SELECT name, creator_id, created_at, id FROM rooms WHERE id=$1`, id)
	if err != nil {
		return nil
	}
	return &room
}

func (repo *Repository) GetRoomParticipants(roomId int64) []int64 {
	var usersId []int64
	err := repo.DB.Select(&usersId, `SELECT user_id FROM rooms_members WHERE room_id=$1`, roomId)
	if err != nil {
		return nil
	}
	return usersId
}

func (repo *Repository) CheckAndGetRoomForUser(roomId, userId int64) *models.Room {
	var room models.Room
	err := repo.DB.Get(&room, `SELECT r.* FROM rooms r
    												 JOIN rooms_members rm ON r.id = rm.room_id 
    												                    	 AND rm.user_id=$1 
    												                    	 AND rm.room_id=$2`, userId, roomId)
	if err != nil {
		return nil
	}
	return &room
}

func (repo *Repository) GetRoomRole(roomId, userId int64) models.RoomRole {
	var role int
	repo.DB.QueryRow(`SELECT role 
													 FROM rooms_members 
													 WHERE room_id=$1 AND user_id=$2`, roomId, userId).Scan(&role)
	return models.RoomRole(role)
}

func (repo *Repository) GetRoomRoles(roomId int64, usersId []int64) []models.RoomMember {
	var roomsMembers []models.RoomMember
	var placeholders []string
	args := []interface{}{roomId}
	for _, userId := range usersId {
		placeholders = append(placeholders, "$"+fmt.Sprint(len(args)+1))
		args = append(args, userId)
	}
	query := fmt.Sprintf("SELECT user_id FROM rooms_members WHERE room_id=$1 AND user_id IN (%s);", strings.Join(placeholders, ", "))
	err := repo.DB.Select(&roomsMembers, query, args...)
	if err != nil {
		return nil
	}
	return roomsMembers
}

func (repo *Repository) GetRoomMember(userId, roomId int64) *models.RoomMember {
	var userRoom models.RoomMember
	err := repo.DB.Get(&userRoom, `SELECT * FROM rooms_members WHERE user_id=$1 AND room_id=$2`, userId, roomId)
	if err != nil {
		return nil
	}
	return &userRoom
}
