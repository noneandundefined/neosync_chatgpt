package smartcaptcha

type ValidateResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Host    string `json:"host,omitempty"`
}
