package token

import "time"

type Maker interface {
	CreateAccessToken(email string, isAdmin bool, duration time.Duration) (*Payload, error)
	VerifyAccessToken(token string) (*Payload, error)

	CreateRefreshToken(email string, isAdmin bool, duration time.Duration) (*Payload, error)
	VerifyRefreshToken(token string) (*Payload, error)
}
