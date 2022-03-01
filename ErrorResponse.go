package auth0

type ErrorResponse struct {
	StatusCode int64  `json:"statusCode"`
	Error      string `json:"error"`
	Message    string `json:"message"`
	ErrorCode  string `json:"errorCode"`
}
