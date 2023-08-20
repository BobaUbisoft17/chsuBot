package database

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func (s *Storage) CreateUsersDatabase() {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
	}
	defer db.Close()
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS users (userID BIGINT, groupID BIGINT)")
	if err != nil {
		s.logger.Error(err)
	}
	defer statement.Close()
	statement.Exec()
}

func (s *Storage) IsUserInDB(userID int64) bool {
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

func (s *Storage) AddUser(userID int64) {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
	}
	defer db.Close()
	InDB := s.IsUserInDB(userID)
	if !InDB {
		_, err := db.Exec("INSERT INTO users (userID, groupID) VALUES ($1)", userID)
		if err != nil {
			s.logger.Error(err)
		}
	}
}

func (s *Storage) IsUserHasGroup(userID int64) bool {
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

func (s *Storage) ChangeUserGroup(userID int64, groupName string) {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
	}
	defer db.Close()
	groupID := s.GroupId(groupName)
	_, err = db.Exec("UPDATE users SET groupID=$1 WHERE userID=$2", groupID, userID)
	if err != nil {
		s.logger.Error(err)
	}
}

func (s *Storage) GetUserGroup(userID int64) int {
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

func (s *Storage) DeleteGroup(userID int64) {
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
