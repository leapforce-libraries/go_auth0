package auth0

import (
	"fmt"
	"net/http"

	o_types "github.com/leapforce-libraries/go_auth0/types"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type User struct {
	CreatedAt     o_types.DateTimeString  `json:"created_at"`
	Email         string                  `json:"email"`
	EmailVerified bool                    `json:"email_verified"`
	Identities    []UserIdentity          `json:"identities"`
	Name          string                  `json:"name"`
	Nickname      string                  `json:"nickname"`
	Picture       string                  `json:"picture"`
	UpdatedAt     o_types.DateTimeString  `json:"updated_at"`
	UserId        string                  `json:"user_id"`
	UserMetadata  map[string]interface{}  `json:"user_metadata"`
	AppMetadata   map[string]interface{}  `json:"app_metadata"`
	LastIp        string                  `json:"last_ip"`
	LastLogin     *o_types.DateTimeString `json:"last_login"`
	LoginsCount   int64                   `json:"logins_count"`
}

type UserIdentity struct {
	Connection string `json:"connection"`
	Provider   string `json:"provider"`
	UserId     string `json:"user_id"`
	IsSocial   bool   `json:"isSocial"`
}

func (service *Service) GetUser(userId string) (*User, *errortools.Error) {
	user := User{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.url(fmt.Sprintf("users/%s", userId)),
		ResponseModel: &user,
	}
	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &user, nil
}
