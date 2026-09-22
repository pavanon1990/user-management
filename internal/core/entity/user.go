package entity

import "time"

type User struct {
	ID           string    `bson:"_id,omitempty"`
	Name         string    `bson:"name"`
	Email        string    `bson:"email"`
	PasswordHash string    `bson:"password"`
	CreatedAt    time.Time `bson:"created_at"`
}
