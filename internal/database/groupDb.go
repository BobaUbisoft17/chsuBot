package database

import (
	"database/sql"

	chsuAPI "github.com/BobaUbisoft17/chsuBot/internal/chsuAPI"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	DbUrl  string
	logger logger
}

type logger interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Error(args ...interface{})
}

func NewStorage(url string, logger logger) *Storage {
	return &Storage{
		DbUrl:  url,
		logger: logger,
	}

}

func (s *Storage) Start() {
	s.CreateGroupDatabase()
	s.CreateUsersDatabase()
	s.logger.Info("Таблицы созданы")
}

func (s *Storage) CreateGroupDatabase() {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
		return
	}
	defer db.Close()
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS groups (groupName TEXT, groupID BIGINT, todaySchedule TEXT, tomorrowSchedule TEXT)")

	if err != nil {
		s.logger.Error(err)
	}
}

func (s *Storage) AddGroups(groupIds []chsuAPI.GroupIds) {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
	}
	defer db.Close()
	_, err = db.Exec("DELETE FROM groups")
	if err != nil {
		s.logger.Error(err)
	}
	statement, err := db.Prepare("INSERT INTO groups (groupName, groupID) VALUES ($1, $2)")
	if err != nil {
		s.logger.Error(err)
	}
	for _, group := range groupIds {
		_, err := statement.Exec(group.Title, group.Id)
		if err != nil {
			s.logger.Error(err)
		}
	}
}

func (s *Storage) GroupId(groupName string) int {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
	}
	defer db.Close()
	row := db.QueryRow("SELECT groupID FROM groups WHERE groupName=$1", groupName)
	var group int
	err = row.Scan(&group)
	if err != nil {
		s.logger.Error(err)
	}
	return group
}

func (s *Storage) GetGroupIds() []int {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT groupID FROM groups")
	if err != nil {
		s.logger.Error(err)
	}
	defer rows.Close()
	groupIds := make([]int, 0)
	var groupId int
	for rows.Next() {
		if err = rows.Scan(&groupId); err != nil {
			s.logger.Error(err)
		}
		groupIds = append(groupIds, groupId)
	}
	return groupIds
}

func (s *Storage) GroupNameIsCorrect(groupName string) bool {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
	}
	defer db.Close()
	row := db.QueryRow("SELECT EXISTS (SELECT groupID FROM groups WHERE groupName=$1 AND groupID IS NOT NULL)", groupName)
	var ans bool
	err = row.Scan(&ans)
	if err != nil {
		s.logger.Error(err)
	}
	return ans
}

func (s *Storage) GetGroupNames() []string {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT groupName FROM groups")
	if err != nil {
		s.logger.Error(err)
	}
	defer rows.Close()
	groupNames := make([]string, 0)
	var groupName string
	for rows.Next() {
		if err = rows.Scan(&groupName); err != nil {
			s.logger.Error(err)
		}
		groupNames = append(groupNames, groupName)
	}
	return groupNames
}

func (s *Storage) UpdateSchedule(todaySchedule, tomorrowSchedule string, groupID int) {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
	}
	defer db.Close()
	_, err = db.Exec("UPDATE groups SET todaySchedule=$1, tomorrowSchedule=$2 WHERE groupID=$3", todaySchedule, tomorrowSchedule, groupID)
	if err != nil {
		s.logger.Error(err)
	}
}

func (s *Storage) UnusedID(IDs []int) []int {
	unusedKeys := []int{}
	usedKeys := make(map[int]bool)
	for _, ID := range IDs {
		usedKeys[ID] = true
	}
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT groupID FROM groups")
	if err != nil {
		s.logger.Error(err)
	}
	defer rows.Close()
	var ID int
	for rows.Next() {
		rows.Scan(&ID)
		if _, ok := usedKeys[ID]; !ok {
			unusedKeys = append(unusedKeys, ID)
		}
	}
	return unusedKeys
}

func (s *Storage) GetTodaySchedule(groupID int) string {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
	}
	defer db.Close()

	row := db.QueryRow("SELECT todaySchedule FROM groups WHERE groupID=$1", groupID)
	var ans string
	err = row.Scan(&ans)
	if err != nil {
		s.logger.Error(err)
	}
	return ans
}

func (s *Storage) GetTomorrowSchedule(groupID int) string {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
	}
	defer db.Close()

	row := db.QueryRow("SELECT tomorrowSchedule FROM groups WHERE groupID=$1", groupID)
	var ans string
	err = row.Scan(&ans)
	if err != nil {
		s.logger.Error(err)
	}
	return ans
}