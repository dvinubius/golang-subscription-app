package data

type DBUsersInterface interface {
	GetAll() ([]*User, error)
	GetByEmail(email string) (*User, error)
	GetOne(id int) (*User, error)
	Update(*User) error
	Delete(*User) error
	Insert(*User) (int, error)
	ResetPassword(*User, string) error
	PasswordMatches(*User, string) (bool, error)
}

type DBPlansInterface interface {
	GetAll() ([]*Plan, error)
	GetOne(id int) (*Plan, error)
	SubscribeUserToPlan(user User, plan Plan) error
	AmountForDisplay(*Plan) string
}
