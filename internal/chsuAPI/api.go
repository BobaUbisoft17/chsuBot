package chsuapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/BobaUbisoft17/chsuBot/internal/schedule"
	"github.com/BobaUbisoft17/chsuBot/pkg/logging"
)

const URL = "http://api.chsu.ru/api/"

var ErrInvalidToken = errors.New("Invalid token")

func New(data map[string]string, logger *logging.Logger) *API {
	return &API{
		Data:   data,
		Client: &http.Client{},
		logger: logger,
	}
}

func (a *API) updateToken() error {
	bytesData, err := json.Marshal(a.Data)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		URL+"/auth/signin",
		"application/json",
		bytes.NewBuffer(bytesData),
	)
	if err != nil {
		return err
	}

	bytesSlice, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	accesToken, err := readJson[Token](bytesSlice)
	if err != nil {
		return err
	}

	a.Token = accesToken.Data
	return nil
}

func (a *API) GroupsId() ([]GroupIds, error) {
	for {
		groupIds, err := a.requestGroupsId()
		if errors.Is(err, ErrInvalidToken) {
			if err = a.updateToken(); err != nil {
				return nil, err
			} else {
				continue
			}
		}
		return groupIds, nil
	}
}

func (a *API) requestGroupsId() ([]GroupIds, error) {
	request, err := http.NewRequest(http.MethodGet, URL+"group/v1", nil)
	if err != nil {
		return []GroupIds{}, err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
	response, err := a.Client.Do(request)
	if err != nil {
		if response.StatusCode == 403 {
			return nil, ErrInvalidToken
		}
		return []GroupIds{}, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return []GroupIds{}, err
	}
	ids, err := readJson[[]GroupIds](body)
	if err != nil {
		if strings.Contains(err.Error(), "invalid character") {
			return nil, ErrInvalidToken
		}
		return []GroupIds{}, err
	}

	return ids, nil
}

func (a *API) One(startDate, endDate string, groupId int) ([]schedule.Lecture, error) {
	for {
		lectures, err := a.requestOne(startDate, endDate, groupId)
		if errors.Is(err, ErrInvalidToken) {
			if err = a.updateToken(); err != nil {
				return nil, err
			} else {
				continue
			}
		}
		return lectures, err
	}
}

func (a *API) requestOne(startDate, endDate string, groupId int) ([]schedule.Lecture, error) {
	requestBody := fmt.Sprintf(
		"timetable/v1/from/%v/to/%v/groupId/%v/",
		startDate,
		endDate,
		groupId,
	)
	request, err := http.NewRequest(http.MethodGet, URL+requestBody, nil)
	if err != nil {
		return []schedule.Lecture{}, err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
	response, err := a.Client.Do(request)
	if err != nil {
		if response.StatusCode == 403 {
			return nil, ErrInvalidToken
		}
		return []schedule.Lecture{}, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return []schedule.Lecture{}, err
	}
	lessons, err := readJson[[]schedule.Lecture](body)
	if err != nil {
		if strings.Contains(err.Error(), "invalid character") {
			return nil, ErrInvalidToken
		}
		return []schedule.Lecture{}, err
	}
	return lessons, nil
}

func (a *API) All() ([]schedule.Lecture, error) {
	for {
		schedules, err := a.requestAll()
		if errors.Is(err, ErrInvalidToken) {
			if err = a.updateToken(); err != nil {
				return nil, err
			} else {
				continue
			}
		}
		return schedules, err
	}
}

func (a *API) requestAll() ([]schedule.Lecture, error) {
	from := time.Now().Format("02.01.2006")
	to := time.Now().Add(24 * time.Hour).Format("02.01.2006")
	requestBody := fmt.Sprintf("timetable/v1/event/from/%s/to/%s/", from, to)
	request, err := http.NewRequest(http.MethodGet, URL+requestBody, nil)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
	if err != nil {
		return []schedule.Lecture{}, err
	}

	resp, err := a.Client.Do(request)
	if err != nil {
		if resp.StatusCode == 403 {
			return nil, ErrInvalidToken
		}
		return []schedule.Lecture{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []schedule.Lecture{}, err
	}
	sliceLectures, err := readJson[[]schedule.Lecture](body)
	if err != nil {
		if strings.Contains(err.Error(), "invalid character") {
			return nil, ErrInvalidToken
		}
		return []schedule.Lecture{}, err
	}
	return sliceLectures, nil
}
