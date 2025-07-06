package token

import "time"

type Maker interface {
	CreateToken(email string, duration time.Duration) (*Payload, error)
	VerifyToken(token string) (*Payload, error)
}
