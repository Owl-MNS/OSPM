package models

type APIError struct {
	Error   string `json:"error"`   // This field determines the raised error
	Message string `json:"message"` // This field determines the detailed information about raise error
}
