package user

type User struct{
	Id	int `json:"-"`
	Username string `json:"name"`
	Password string `json:"password"`
}