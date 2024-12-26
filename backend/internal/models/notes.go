package models

type Note struct {
	Id    string `json:"Id,omitempty"`
	Text  string `json:"text,omitempty"`
	Token string `json:"captcha_token,omitempty"`
}
