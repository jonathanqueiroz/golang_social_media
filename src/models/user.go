package models

import (
	"errors"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint64    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Birthdate string    `json:"birthdate,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Password  string    `json:"password,omitempty"`
}

func (user *User) Prepare(method string) error {
	if err := user.validate(method); err != nil {
		return err
	}

	if method == "create" {
		if err := user.format(method); err != nil {
			return err
		}
	}

	return nil
}

func (user *User) validate(method string) error {
	if user.Name == "" {
		return errors.New("name is required")
	}

	if user.Email == "" {
		return errors.New("email is required")
	}

	if err := checkmail.ValidateFormat(user.Email); err != nil {
		return errors.New("invalid email format")
	}

	if user.Birthdate == "" {
		return errors.New("birthdate is required")
	}

	_, err := time.Parse("01-02-2006", user.Birthdate)
	if err != nil {
		return errors.New("invalid birthdate format, use DD-MM-YYYY")
	}

	if method == "create" && user.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

func (user *User) format(method string) error {
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(user.Email)
	user.Birthdate = strings.TrimSpace(user.Birthdate)

	if method == "create" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		user.Password = string(hashedPassword)
	}

	parsedBirthdate, _ := time.Parse("01-02-2006", user.Birthdate)
	user.Birthdate = parsedBirthdate.Format("2006-01-02")

	return nil
}
