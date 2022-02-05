package storage

type User struct {
	UserName string `json:"userName"`
	Name string `json:"name"`
	LastName string `json:"lastName"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type UserStorage interface {
	Init()
	Get(userName string) (User, error)
	GetUsers(query string) ([]User, error)
	Add(usr User) (error)
	Update(usr User) (error)
	Delete(userName string) (error)
	SetPassword(userName string, password string)
}
