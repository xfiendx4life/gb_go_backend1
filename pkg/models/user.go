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
