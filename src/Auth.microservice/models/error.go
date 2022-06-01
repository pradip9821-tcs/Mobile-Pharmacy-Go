package models

type Error struct {
	Message          string `json:"message"`
	Code             string `json:"code"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	Status           int    `json:"status"`
}
