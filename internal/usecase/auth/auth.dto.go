package auth

import "github.com/golang-jwt/jwt/v5"

type SignIn struct {
	Email string `json:"email"`
}

type SignUp struct {
	Email string `json:"email"`
}
type SendEmailCode struct {
	Email string
}

type GenerateToken struct {
	Id   int64
	Role int
}

type CheckCode struct {
	Token string `json:"token"`
	Code  string `json:"code"`
}

type ResendData struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type UpdatePsw struct {
	Token    string
	Password string
}

type ResendCode struct {
	Token string
}

type Claims struct {
	ID    int
	Email string
	Role  int
	jwt.RegisteredClaims
}
