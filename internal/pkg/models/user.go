package models

type Userer interface {
	Save() (id int, err error)
	Get(name, password string)
}

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func NewUser(name, password, email string) *User {
	return &User{
		Name:     name,
		Password: password,
		Email:    email,
	}
}
