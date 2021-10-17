package v1

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"error"`
}

type ListMetaData struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Count  int `json:"count"`
}
