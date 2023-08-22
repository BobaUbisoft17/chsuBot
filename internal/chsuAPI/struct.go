package chsuapi

import (
	"net/http"

	"github.com/BobaUbisoft17/chsuBot/pkg/logging"
)

type API struct {
	Data   map[string]string
	Token  string
	Client *http.Client
	logger *logging.Logger
}

type Token struct {
	Data  string `json:"data"`
	Error struct {
		Code        int    `json:"code"`
		Status      string `json:"status"`
		Description string `json:"description"`
	}
}

type GroupIds struct {
	Id    int
	Title string
}
