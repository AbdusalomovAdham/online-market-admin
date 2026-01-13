package auth

import "main/internal/entity"

type SignIn struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type GetInfo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Token     string `json:"token"`
}

type AdminDetails struct {
	Id        int64        `json:"id"`
	Role      int          `json:"role"`
	Password  string       `json:"password"`
	FirstName string       `json:"first_name"`
	LastName  string       `json:"last_name"`
	Avatar    *entity.File `json:"avatar"`
}
