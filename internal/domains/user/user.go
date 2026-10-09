package user

import "uuid"

type User struct {
	ID       uuid.UUID
	FullName string
}
