package logic

import (
	"errors"
	"fmt"
	"regexp"
)

const (
	minLoginLen  = 4
	maxLoginLen  = 32
	loginPattern = "^[a-zA-Z0-9_\\-]+$"

	minPasswordLen  = 8
	maxPasswordLen  = 64
	passwordPattern = "^[a-zA-Z0-9!@#\\$%\\^&\\*()_+=\\-]+$"
)

var (
	loginRegexp    = regexp.MustCompile(loginPattern)
	passwordRegexp = regexp.MustCompile(passwordPattern)
)

func validateLogin(login string) error {
	n := len(login)
	if n < minLoginLen || n > maxLoginLen {
		return fmt.Errorf("логин должен содержать от %d до %d символов", minLoginLen, maxLoginLen)
	}
	if !loginRegexp.MatchString(login) {
		return errors.New("логин может содержать только латинские буквы, цифры, точку и дефис")
	}
	return nil
}

func validatePassword(password string) error {
	n := len(password)
	if n < minPasswordLen || n > maxPasswordLen {
		return fmt.Errorf("пароль должен содержать от %d до %d символов", minPasswordLen, maxPasswordLen)
	}
	if !passwordRegexp.MatchString(password) {
		return errors.New("пароль может содержать только буквы, цифры и символы !@#$%^&*()_+=-")
	}
	return nil
}
