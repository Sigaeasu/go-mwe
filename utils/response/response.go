package response

type ResponseAPI struct {
	Status string `json:"status,omitempty"`
	Data interface{} `json:"data,omitempty"`
	Error_ *ApiError `json:"error,omitempty"`
}

type ApiError struct {
	Error string `json: "error`
}
