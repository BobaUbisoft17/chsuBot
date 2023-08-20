package parse

import "net/http"

type logger interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}

type API struct {
	Data   map[string]string
	Token  string
	Client *http.Client
	logger logger
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
