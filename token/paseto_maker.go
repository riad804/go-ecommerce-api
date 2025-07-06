package token

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	bytes, err := base64.StdEncoding.DecodeString(symmetricKey)
	if err != nil {
		return nil, err
	}
	if len(bytes) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid symmetric key size %d: must be %d bytes", len(symmetricKey), chacha20poly1305.KeySize)
	}
	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}
	return maker, nil
}

func (maker *PasetoMaker) CreateToken(email string, duration time.Duration) (*Payload, error) {
	payload, err := NewPayload(email, duration)
	if err != nil {
		return nil, err
	}

	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, err
	}

	payload.Token = token
	return payload, nil
}
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if err := payload.Valid(); err != nil {
		return nil, err
	}
	payload.Token = token

	return payload, nil
}
