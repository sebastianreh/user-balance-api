package user

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	IsDeleted bool   `json:"-"`
}

func CreateUserByRecord(record []string) User {
	return User{
		ID: record[1],
	}
}
