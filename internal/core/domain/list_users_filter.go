package domain

import "time"

type ListUsersFilter struct {
	Limit  int64
	Offset int64
	From   time.Time
	To     time.Time
}
