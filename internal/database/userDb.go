package database

import (
	"database/sql"

	"github.com/BobaUbisoft17/chsuBot/pkg/logging"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type UserStorage struct {
	DbUrl  string
	gs     groupStorage
	logger *logging.Logger
}

type groupStorage interface {
	GroupId(string) int
}

func NewUserStorage(url string, gs groupStorage, logger *logging.Logger) *UserStorage {
	return &UserStorage{
		DbUrl:  url,
		gs:     gs,
		logger: logger,
	}
}

func (s *UserStorage) Start() {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
	}
	defer db.Close()
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (userID BIGINT, groupID BIGINT)")
	if err != nil {
		s.logger.Error(err)
	}
}

func (s *UserStorage) IsUserInDB(userID int64) bool {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
	}
	defer db.Close()
	row := db.QueryRow("SELECT EXISTS (SELECT userID FROM users WHERE userID=$1)", userID)
	var ans bool
	err = row.Scan(&ans)
	if err != nil {
		s.logger.Error(err)
	}
	return ans
}

func (s *UserStorage) AddUser(userID int64) {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
	}
	defer db.Close()
	InDB := s.IsUserInDB(userID)
	if !InDB {
		_, err := db.Exec("INSERT INTO users (userID, groupID) VALUES ($1, NULL)", userID)
		if err != nil {
			s.logger.Error(err)
		}
	}
}

func (s *UserStorage) IsUserHasGroup(userID int64) bool {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
	}
	defer db.Close()
	row := db.QueryRow("SELECT EXISTS (SELECT FROM users WHERE userID=$1 AND groupID IS NOT NULL)", userID)
	var ans bool
	if err := row.Scan(&ans); err != nil {
		s.logger.Error(err)
	}
	return ans
}

func (s *UserStorage) ChangeUserGroup(userID int64, groupID int) {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
	}
	defer db.Close()
	_, err = db.Exec("UPDATE users SET groupID=$1 WHERE userID=$2", groupID, userID)
	if err != nil {
		s.logger.Error(err)
	}
}

func (s *UserStorage) GetUserGroup(userID int64) int {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
	}
	defer db.Close()
	row := db.QueryRow("SELECT groupID FROM users WHERE userID=$1", userID)
	var group int
	err = row.Scan(&group)
	if err != nil {
		s.logger.Error(err)
	}
	return group
}

func (s *UserStorage) DeleteGroup(userID int64) {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
	}
	defer db.Close()
	_, err = db.Exec("UPDATE users SET groupID=NULL WHERE userID=$1", userID)
	if err != nil {
		s.logger.Error(err)
	}
}

func (s *UserStorage) DeleteUser(userID int64) {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
	}
	defer db.Close()
	_, err = db.Exec("DELETE FROM users WHERE userID=$1", userID)
	if err != nil {
		s.logger.Error(err)
	}
}

func (s *UserStorage) GetUsersId() []int {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Errorf("%v", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT userID FROM users")
	if err != nil {
		s.logger.Errorf("%v", err)
	}
	defer rows.Close()

	var userIds []int
	var id int
	for rows.Next() {
		if err = rows.Scan(&id); err != nil {
			s.logger.Errorf("%v", err)
		}
		userIds = append(userIds, id)
	}
	return userIds
}
