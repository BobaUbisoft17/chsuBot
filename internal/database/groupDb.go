package database

import (
	"database/sql"
	"slices"

	chsuAPI "github.com/BobaUbisoft17/chsuBot/internal/chsuAPI"
	"github.com/BobaUbisoft17/chsuBot/internal/schedule"
	"github.com/BobaUbisoft17/chsuBot/pkg/logging"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type GroupStorage struct {
	DbUrl  string
	logger *logging.Logger
}

type GroupInfo struct {
	GroupName string
	GroupID   int
}

func NewGroupStorage(url string, logger *logging.Logger) *GroupStorage {
	return &GroupStorage{
		DbUrl:  url,
		logger: logger,
	}

}

func (s *GroupStorage) Start() error {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	defer db.Close()
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS groups (groupName TEXT, groupID BIGINT, todaySchedule TEXT, tomorrowSchedule TEXT)")
	if err != nil {
		s.logger.Error(err)
		return err
	}
	return nil
}

func (s *GroupStorage) AddGroups(groupIds []chsuAPI.GroupIds) error {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	defer db.Close()
	_, err = db.Exec("DELETE FROM groups")
	if err != nil {
		s.logger.Error(err)
		return err
	}
	statement, err := db.Prepare("INSERT INTO groups (groupName, groupID) VALUES ($1, $2)")
	if err != nil {
		s.logger.Error(err)
		return err
	}
	for _, group := range groupIds {
		_, err := statement.Exec(group.Title, group.Id)
		if err != nil {
			s.logger.Error(err)
			return err
		}
	}
	return nil
}

func (s *GroupStorage) GetGroupIds() ([]int, error) {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query("SELECT groupID FROM groups")
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	defer rows.Close()
	groupIds := make([]int, 0)
	var groupId int
	for rows.Next() {
		if err = rows.Scan(&groupId); err != nil {
			s.logger.Error(err)
			return nil, err
		}
		groupIds = append(groupIds, groupId)
	}
	return groupIds, nil
}

func (s *GroupStorage) GroupNameIsCorrect(groupName string) (bool, error) {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
	}
	defer db.Close()
	row := db.QueryRow("SELECT EXISTS (SELECT groupID FROM groups WHERE groupName=$1 AND groupID IS NOT NULL)", groupName)
	var ans bool
	if err = row.Scan(&ans); err != nil {
		s.logger.Error(err)
		return false, err
	}
	return ans, nil
}

func (s *GroupStorage) GetGroupNames() ([]string, error) {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query("SELECT groupName FROM groups")
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	defer rows.Close()
	groupNames := make([]string, 0)
	var groupName string
	for rows.Next() {
		if err = rows.Scan(&groupName); err != nil {
			s.logger.Error(err)
			return nil, err
		}
		groupNames = append(groupNames, groupName)
	}
	return groupNames, nil
}

func (s *GroupStorage) UpdateSchedule(todaySchedule, tomorrowSchedule []schedule.Lecture, groupID int) error {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	defer db.Close()
	_, err = db.Exec(
		"UPDATE groups SET todaySchedule=$1, tomorrowSchedule=$2 WHERE groupID=$3",
		schedule.New(todaySchedule).Render(),
		schedule.New(tomorrowSchedule).Render(),
		groupID,
	)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	return nil
}

func (s *GroupStorage) UnusedID(ids []int) ([]int, error) {
	unusedKeys := []int{}
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query("SELECT groupID FROM groups")
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	defer rows.Close()
	var id int
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			s.logger.Errorf("%v", err)
			return nil, err
		}
		if !slices.Contains(ids, id) {
			unusedKeys = append(unusedKeys, id)
		}
	}
	return unusedKeys, nil
}

func (s *GroupStorage) GetTodaySchedule(groupID int) (string, error) {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
		return "", err
	}
	defer db.Close()

	row := db.QueryRow("SELECT todaySchedule FROM groups WHERE groupID=$1", groupID)
	var ans string
	if err = row.Scan(&ans); err != nil {
		s.logger.Error(err)
		return "", err
	}
	return ans, nil
}

func (s *GroupStorage) GetTomorrowSchedule(groupID int) (string, error) {
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
		return "", err
	}
	defer db.Close()

	row := db.QueryRow("SELECT tomorrowSchedule FROM groups WHERE groupID=$1", groupID)
	var ans string
	if err = row.Scan(&ans); err != nil {
		s.logger.Error(err)
		return "", err
	}
	return ans, nil
}

func (s *GroupStorage) GroupsStartsWith(firstSymbol string) ([]GroupInfo, error) {
	var groups []GroupInfo
	db, err := sql.Open("pgx", s.DbUrl)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT groupName, groupID FROM groups WHERE groupName LIKE $1||'%' ORDER BY groupName", firstSymbol)
	if err != nil {
		s.logger.Errorf("%v", err)
		return nil, err
	}
	defer rows.Close()

	var group GroupInfo
	for rows.Next() {
		if err = rows.Scan(&group.GroupName, &group.GroupID); err != nil {
			s.logger.Errorf("%v", err)
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}
