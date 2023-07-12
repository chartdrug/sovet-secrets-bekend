package entity

import "time"

type UsersQuestion struct {
	Id        int       `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Username  string    `json:"username"`
	ChatId    int       `json:"chatId"`
	Message   string    `json:"message"`
	Answer    string    `json:"answer"`
	AnswerTxt string    `json:"answer_txt"`
}

// GetID returns the user ID.
func (u UsersQuestion) GetID() int {
	return u.Id
}
