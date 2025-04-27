package controller

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *apiError) Error() string {
	return e.Message
}
