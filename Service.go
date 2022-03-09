package auth0

import (
	"fmt"
	"net/http"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	go_oauth2 "github.com/leapforce-libraries/go_oauth2"
	tokensource "github.com/leapforce-libraries/go_oauth2/tokensource"
	utilities "github.com/leapforce-libraries/go_utilities"
)

const (
	apiName string = "Auth0"
)

// type
//
type Service struct {
	tenantName    string
	clientId      string
	oAuth2Service *go_oauth2.Service
}

type ServiceConfig struct {
	TenantName  string
	ClientId    string
	TokenSource tokensource.TokenSource
}

func NewService(serviceConfig *ServiceConfig) (*Service, *errortools.Error) {
	if serviceConfig == nil {
		return nil, errortools.ErrorMessage("ServiceConfig must not be a nil pointer")
	}

	if serviceConfig.TenantName == "" {
		return nil, errortools.ErrorMessage("Service TenantName not provided")
	}

	if serviceConfig.ClientId == "" {
		return nil, errortools.ErrorMessage("Service ClientId not provided")
	}

	config := go_oauth2.ServiceConfig{
		ClientId:    serviceConfig.ClientId,
		TokenSource: serviceConfig.TokenSource,
	}

	oAuth2Service, e := go_oauth2.NewService(&config)
	if e != nil {
		return nil, e
	}

	return &Service{
		tenantName:    serviceConfig.TenantName,
		clientId:      serviceConfig.ClientId,
		oAuth2Service: oAuth2Service,
	}, nil
}

func (service *Service) httpRequest(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	// add authentication header
	errorResponse := ErrorResponse{}
	if utilities.IsNil(requestConfig.ErrorModel) {
		// add error model
		(*requestConfig).ErrorModel = &errorResponse
	}

	request, response, e := service.oAuth2Service.HttpRequest(requestConfig)
	if e != nil {
		if errorResponse.Message != "" {
			e.SetMessage(errorResponse.Message)
		}
	}

	return request, response, e
}

func (service *Service) url(path string) string {
	return fmt.Sprintf("https://%s.us.auth0.com/api/v2/%s", service.tenantName, path)
}

func (service *Service) ApiName() string {
	return apiName
}

func (service *Service) ApiKey() string {
	return service.clientId
}

func (service *Service) ApiCallCount() int64 {
	return service.oAuth2Service.ApiCallCount()
}

func (service *Service) ApiReset() {
	service.oAuth2Service.ApiReset()
}
