package domain

// User is a domain User.
type User struct {
	ID       int
	Username string
	Password string
	Admin    bool
}

type NewUserData struct {
	ID       int
	Username string
	Password string
	Admin    bool
}

// NewUser creates a new user.
func NewUser(data NewUserData) (User, error) {
	return User(data), nil
}
