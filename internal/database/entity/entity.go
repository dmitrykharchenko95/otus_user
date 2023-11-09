package entity

import (
	"fmt"
	"strings"
)

type User struct {
	Id        int64
	Username  string
	FirstName string
	LastName  string
	Email     string
	Phone     string
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
