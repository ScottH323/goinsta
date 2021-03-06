package goinsta

import (
	"net/http"
	"net/http/cookiejar"
	"github.com/scotth323/goinsta/response"
)

type Informations struct {
	Username  string
	Password  string
	DeviceID  string
	UUID      string
	RankToken string
	Token     string
	PhoneID   string
}

type Instagram struct {
	Cookiejar *cookiejar.Jar
	InstaType
	Transport http.Transport
}

type InstaType struct {
	IsLoggedIn   bool
	Informations Informations
	LoggedInUser response.User

	Proxy string
}

type BackupType struct {
	Cookies []http.Cookie
	InstaType
}
