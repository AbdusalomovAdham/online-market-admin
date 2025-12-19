package auth

type SignIn struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type GetInfo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Token     string `json:"token"`
}
