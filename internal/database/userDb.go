package database

import (
	"database/sql"

	"github.com/BobaUbisoft17/chsuBot/pkg/logging"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type UserStorage struct {
	DbUrl  string
	logger *logging.Logger
}

func NewUserStorage(url string, logger *logging.Logger) *UserStorage {
	return &UserStorage{
		DbUrl:  url,
		logger: logger,
	}
}

func (s *UserStorage) Start() error {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	defer db.Close()
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (userID BIGINT, groupID BIGINT)")
	if err != nil {
		s.logger.Error(err)
		return err
	}
	return nil
}

func (s *UserStorage) IsUserInDB(userID int64) (bool, error) {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
		return false, err
	}
	defer db.Close()
	row := db.QueryRow("SELECT EXISTS (SELECT userID FROM users WHERE userID=$1)", userID)
	var ans bool
	if err = row.Scan(&ans); err != nil {
		s.logger.Error(err)
		return false, err
	}
	return ans, nil
}

func (s *UserStorage) AddUser(userID int64) error {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	defer db.Close()
	InDB, err := s.IsUserInDB(userID)
	if err != nil {
		return err
	}
	if !InDB {
		_, err := db.Exec("INSERT INTO users (userID, groupID) VALUES ($1, NULL)", userID)
		if err != nil {
			s.logger.Error(err)
			return err
		}
	}
	return nil
}

func (s *UserStorage) IsUserHasGroup(userID int64) (bool, error) {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
		return false, err
	}
	defer db.Close()
	row := db.QueryRow("SELECT EXISTS (SELECT FROM users WHERE userID=$1 AND groupID IS NOT NULL)", userID)
	var ans bool
	if err := row.Scan(&ans); err != nil {
		s.logger.Error(err)
		return false, err
	}
	return ans, nil
}

func (s *UserStorage) ChangeUserGroup(userID int64, groupID int) error {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	defer db.Close()
	_, err = db.Exec("UPDATE users SET groupID=$1 WHERE userID=$2", groupID, userID)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	return nil
}

func (s *UserStorage) GetUserGroup(userID int64) (int, error) {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
		return 0, err
	}
	defer db.Close()
	row := db.QueryRow("SELECT groupID FROM users WHERE userID=$1", userID)
	var group int
	if err = row.Scan(&group); err != nil {
		s.logger.Error(err)
		return 0, err
	}
	return group, nil
}

func (s *UserStorage) DeleteGroup(userID int64) error {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	defer db.Close()
	_, err = db.Exec("UPDATE users SET groupID=NULL WHERE userID=$1", userID)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	return nil
}

func (s *UserStorage) DeleteUser(userID int64) error {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	defer db.Close()
	_, err = db.Exec("DELETE FROM users WHERE userID=$1", userID)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	return nil
}

func (s *UserStorage) GetUsersId() ([]int, error) {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Errorf("%v", err)
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT userID FROM users")
	if err != nil {
		s.logger.Errorf("%v", err)
		return nil, err
	}
	defer rows.Close()

	var userIds []int
	var id int
	for rows.Next() {
		if err = rows.Scan(&id); err != nil {
			s.logger.Errorf("%v", err)
			return nil, err
		}
		userIds = append(userIds, id)
	}
	return userIds, nil
}
