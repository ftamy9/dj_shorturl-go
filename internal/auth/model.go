package auth

import "time"

type BaseUser struct {
	ID           string    `json:"id"`
	PasswordHash string    `json:"-"`
	CreateDate   time.Time `json:"create_date"`
}
