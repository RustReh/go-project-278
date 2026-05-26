package schemas

type ValidationErrorsResponse struct {
	Errors map[string]string `json:"errors"`
}

type InvalidRequestResponse struct {
	Error string `json:"error"`
}
