package types

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func NewUser(id, username string) *User {
	return &User{
		ID:       id,
		Username: username,
	}
}
