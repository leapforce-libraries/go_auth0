package auth0

import (
	"fmt"
	"net/http"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	utilities "github.com/leapforce-libraries/go_utilities"
)

const (
	apiName string = "Auth0"
)

// type
//
type Service struct {
	apiUrl      string
	apiKey      string
	httpService *go_http.Service
}

type ServiceConfig struct {
	ApiUrl string
	ApiKey string
}

func NewService(serviceConfig *ServiceConfig) (*Service, *errortools.Error) {
	if serviceConfig == nil {
		return nil, errortools.ErrorMessage("ServiceConfig must not be a nil pointer")
	}

	if serviceConfig.ApiUrl == "" {
		return nil, errortools.ErrorMessage("Service ApiKey not provided")
	}

	if serviceConfig.ApiKey == "" {
		return nil, errortools.ErrorMessage("Service ApiKey not provided")
	}

	httpService, e := go_http.NewService(&go_http.ServiceConfig{})
	if e != nil {
		return nil, e
	}

	return &Service{
		apiUrl:      strings.Trim(serviceConfig.ApiUrl, "/"),
		apiKey:      serviceConfig.ApiKey,
		httpService: httpService,
	}, nil
}

func (service *Service) httpRequest(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	// add authentication header
	header := http.Header{}
	header.Set("Authorization", fmt.Sprintf("Bearer %s", service.apiKey))
	(*requestConfig).NonDefaultHeaders = &header

	errorResponse := ErrorResponse{}
	if utilities.IsNil(requestConfig.ErrorModel) {
		// add error model
		(*requestConfig).ErrorModel = &errorResponse
	}

	request, response, e := service.httpService.HttpRequest(requestConfig)
	if e != nil {
		if errorResponse.Message != "" {
			e.SetMessage(errorResponse.Message)
		}
	}

	return request, response, e
}

func (service *Service) url(path string) string {
	return fmt.Sprintf("%s/%s", service.apiUrl, path)
}

func (service *Service) ApiName() string {
	return apiName
}

func (service *Service) ApiKey() string {
	return service.apiKey
}

func (service *Service) ApiCallCount() int64 {
	return service.httpService.RequestCount()
}

func (service *Service) ApiReset() {
	service.httpService.ResetRequestCount()
}
