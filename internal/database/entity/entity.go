package entity

import (
	"errors"
	"fmt"
	"strings"

	"github.com/xyproto/randomstring"
	"golang.org/x/crypto/bcrypt"
)

const localSalt = "D@j~V#2g"

var ErrInvalidPass = errors.New("password invalid")

type User struct {
	Id           int64
	Username     string
	FirstName    string
	LastName     string
	Email        string
	Phone        string
	Salt         string
	PasswordHash string
}

func (u *User) GetUpdateQuery() (setQuery string, args []interface{}) {
	if u.Id == 0 {
		return "", nil
	}

	var (
		fields  []string
		counter int
	)

	if u.FirstName != "" {
		fields = append(fields, "first_name")
		args = append(args, u.FirstName)
		counter++
	}

	if u.LastName != "" {
		fields = append(fields, "last_name")
		args = append(args, u.LastName)
		counter++
	}

	if u.Email != "" {
		fields = append(fields, "email")
		args = append(args, u.Email)
		counter++
	}

	if u.Phone != "" {
		fields = append(fields, "phone")
		args = append(args, u.Phone)
		counter++
	}

	var peaces []string
	for i := 1; i <= counter; i++ {
		peaces = append(peaces, fmt.Sprintf("%s = $%d", fields[i-1], i))
	}

	return fmt.Sprintf("SET %s, updated_at = NOW() WHERE id = %d", strings.Join(peaces, ", "), u.Id), args
}

func (u *User) SetPassword(password string) error {
	var salt = randomstring.HumanFriendlyString(10)

	saltedPassword := password + localSalt + salt
	bHash, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Salt = salt
	u.PasswordHash = string(bHash)
	return nil
}

func (u *User) ValidatePassword(password string) error {
	saltedPassword := password + localSalt + u.Salt
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(saltedPassword)) != nil {
		return ErrInvalidPass
	}
	return nil
}
