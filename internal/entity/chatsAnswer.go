package entity

type ChatsAnswer struct {
	Id       int    `json:"id"`
	Answer   string `json:"answer"`
	Category string `json:"category"`
}

// GetID returns the user ID.
func (u ChatsAnswer) GetID() int {
	return u.Id
}
