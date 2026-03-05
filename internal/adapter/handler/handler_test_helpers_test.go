package handler_test

// errorBody represents the standard JSON error envelope returned by ErrorResponse.
type errorBody struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}
