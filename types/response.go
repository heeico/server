package types

type SuccessResponse struct {
	Message string                 `json:"message"`
	Status  bool                   `json:"status"`
	Data    map[string]interface{} `json:"data"`
}

type FailResponse struct {
	Error  string                 `json:"error"`
	Status bool                   `json:"status"`
	Data   map[string]interface{} `json:"data,omitempty"`
}

type ResponseData map[string]interface{}
