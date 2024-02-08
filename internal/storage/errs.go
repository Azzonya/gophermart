package storage

import (
	"fmt"
)

type ErrUserNotUniq struct {
	Login string
}

func (err ErrUserNotUniq) Error() string {
	return fmt.Sprintf("user with login \"%s\" already exists", err.Login)
}
