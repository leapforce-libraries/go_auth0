package auth0

import (
	"fmt"
	"net/http"

	errortools "github.com/leapforce-libraries/go_errortools"
	gcs "github.com/leapforce-libraries/go_googlecloudstorage"
	go_http "github.com/leapforce-libraries/go_http"
	token "github.com/leapforce-libraries/go_oauth2/token"
)

type TokenMap struct {
	token        *token.Token
	map_         *gcs.Map
	tenantName   string
	clientId     string
	clientSecret string
	audience     string
	httpService  *go_http.Service
}

func NewTokenMap(map_ *gcs.Map, tenantName string, clientId string, clientSecret string, audience string) (*TokenMap, *errortools.Error) {
	if map_ == nil {
		return nil, errortools.ErrorMessage("Map is a nil pointer")
	}

	httpService, e := go_http.NewService(&go_http.ServiceConfig{})
	if e != nil {
		return nil, e
	}

	return &TokenMap{
		map_:         map_,
		tenantName:   tenantName,
		clientId:     clientId,
		clientSecret: clientSecret,
		audience:     audience,
		httpService:  httpService,
	}, nil
}

func (m *TokenMap) Token() *token.Token {
	return m.token
}

func (m *TokenMap) NewToken() (*token.Token, *errortools.Error) {
	t := token.Token{}

	body := struct {
		GrantType    string `json:"grant_type"`
		ClientId     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Audience     string `json:"audience"`
	}{
		"client_credentials",
		m.clientId,
		m.clientSecret,
		m.audience,
	}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           fmt.Sprintf("https://%s.us.auth0.com/oauth/token", m.tenantName),
		BodyModel:     body,
		ResponseModel: &t,
	}
	_, _, e := m.httpService.HttpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &t, nil
}

func (m *TokenMap) SetToken(token *token.Token, save bool) *errortools.Error {
	if token.AccessToken == nil {
		return errortools.ErrorMessage("AccessToken of new token is nil")
	}

	m.token = token

	if !save {
		return nil
	}

	return m.SaveToken()
}

func (m *TokenMap) RetrieveToken() *errortools.Error {
	accessToken, _ := m.map_.Get("access_token")
	tokenType, _ := m.map_.Get("token_type")
	scope, _ := m.map_.Get("scope")
	expiry, _ := m.map_.GetTimestamp("expiry")

	m.token = &token.Token{
		AccessToken: accessToken,
		TokenType:   tokenType,
		Scope:       scope,
		Expiry:      expiry,
	}

	return nil
}

func (m *TokenMap) SaveToken() *errortools.Error {
	if m.token == nil {
		return errortools.ErrorMessage("Token is nil pointer")
	}

	if m.token.AccessToken != nil {
		m.map_.Set("access_token", *m.token.AccessToken, false)
	}

	if m.token.TokenType != nil {
		m.map_.Set("token_type", *m.token.TokenType, false)
	}

	if m.token.Scope != nil {
		m.map_.Set("scope", *m.token.Scope, false)
	}

	if m.token.Expiry != nil {
		m.map_.SetTimestamp("expiry", *m.token.Expiry, false)
	}

	e := m.map_.Save()
	if e != nil {
		return e
	}

	return nil
}
