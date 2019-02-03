package senrenrpc

type SuccessError struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type SuccessErrorOnly struct {
	SuccessError `json:"result"`
}
