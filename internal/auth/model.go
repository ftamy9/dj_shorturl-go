package auth

type BaseUser struct {
	ID           string `json:"id"`
	PasswordHash string `json:"-"`
	CreateDate   string `json:"create_date"`
}
