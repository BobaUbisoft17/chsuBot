package parse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/BobaUbisoft17/chsuBot/internal/schedule"
)

const URL = "http://api.chsu.ru/api/"

func New(data map[string]string, logger logger) *API {
	return &API{
		Data:   data,
		Client: &http.Client{},
		logger: logger,
	}
}

func (a *API) tokenIsValid() (bool, error) {
	if a.Token == "" {
		return false, nil
	}

	resp, err := http.Post(URL+"auth/valid/", "application/json", bytes.NewBufferString(a.Token))
	if err != nil {
		return false, err
	}

	byteResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	answer, err := strconv.ParseBool(string(byteResp))
	if err != nil {
		return false, err
	}

	return answer, nil

}

func (a *API) updateToken() error {
	bytesDate, err := json.Marshal(a.Data)
	if err != nil {
		return err
	}

	resp, err := http.Post(URL+"/auth/signin", "application/json", bytes.NewBuffer(bytesDate))
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
	tokenIsValid, err := a.tokenIsValid()
	if err != nil {
		return []GroupIds{}, err
	}

	if !tokenIsValid {
		if err = a.updateToken(); err != nil {
			return []GroupIds{}, err
		}
	}

	request, err := http.NewRequest(http.MethodGet, URL+"group/v1", nil)
	if err != nil {
		return []GroupIds{}, err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
	response, err := a.Client.Do(request)
	if err != nil {
		return []GroupIds{}, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return []GroupIds{}, err
	}
	ids, err := readJson[[]GroupIds](body)
	if err != nil {
		return []GroupIds{}, err
	}

	return ids, nil
}

func (a *API) One(startDate, endDate string, groupId int) ([]schedule.Lecture, error) {
	request_budy := fmt.Sprintf("timetable/v1/from/%v/to/%v/groupId/%v/", startDate, endDate, groupId)
	request, err := http.NewRequest(http.MethodGet, URL+request_budy, nil)
	if err != nil {
		return []schedule.Lecture{}, err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
	response, err := a.Client.Do(request)
	if err != nil {
		return []schedule.Lecture{}, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return []schedule.Lecture{}, err
	}
	lessons, err := readJson[[]schedule.Lecture](body)
	if err != nil {
		if strings.Contains(err.Error(), "invalid character") {
			if err = a.updateToken(); err != nil {
				return []schedule.Lecture{}, err
			}
			return a.One(startDate, endDate, groupId)
		} else {
			return []schedule.Lecture{}, err
		}
	}
	return lessons, nil
}

func (a *API) All() ([]schedule.Lecture, error) {
	from := time.Now().Format("02.01.2006")
	to := time.Now().Add(24 * time.Hour).Format("02.01.2006")
	request_body := fmt.Sprintf("timetable/v1/event/from/%s/to/%s/", from, to)
	request, err := http.NewRequest(http.MethodGet, URL+request_body, nil)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
	if err != nil {
		return []schedule.Lecture{}, err
	}

	resp, err := a.Client.Do(request)
	if err != nil {
		return []schedule.Lecture{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []schedule.Lecture{}, err
	}
	sliceLectures, err := readJson[[]schedule.Lecture](body)
	if err != nil {
		if strings.Contains(err.Error(), "invalid character") {
			if err = a.updateToken(); err != nil {
				return []schedule.Lecture{}, err
			}
			return a.All()
		} else {
			return []schedule.Lecture{}, err
		}
	}
	return sliceLectures, nil
}