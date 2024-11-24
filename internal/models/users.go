package models

import (
	"kyri56xcaesar/discord_bots_app/internal/utils"
)

type User struct {
	Username string
	Password string
}

func (u *User) VerifyUser() error {
	if !utils.IsValidUsername(u.Username) {
		return &utils.FieldError{Field: "Username", Message: "invalid username"}
	}

	if !utils.IsValidPassword(u.Password) {
		return &utils.FieldError{Field: "Password", Message: "invalid password"}
	}
	return nil
}
