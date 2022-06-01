package models

type Error struct {
	Message          string `json:"message"`
	Error            string `json:"error"`
	Code             string `json:"code"`
	ErrorDescription string `json:"error_description"`
	Status           int    `json:"status"`
}
