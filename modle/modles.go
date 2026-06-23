package modle

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username,omitempty"`
	Age      int    `json:"age"`
	Email    string `json:"email"`
}

var users = []User{
	{
		Id:       1,
		Username: "Kamal",
		Age:      25,
		Email:    "kamal@mail.com",
	},
	{
		Id:       2,
		Username: "Jamal",
		Age:      20,
		Email:    "jamal@mail.com",
	},
}
