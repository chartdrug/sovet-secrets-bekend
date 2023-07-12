package entity

type ChatsQuestion struct {
	Id       int    `json:"id"`
	Question string `json:"question"`
	AnswerId int    `json:"answer_id"`
}

// GetID returns the user ID.
func (u ChatsQuestion) GetID() int {
	return u.Id
}
