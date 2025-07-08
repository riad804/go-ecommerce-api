package token

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

type PasetoMaker struct {
	paseto     *paseto.V2
	AccessKey  []byte
	RefreshKey []byte
}

func NewPasetoMaker(accessKey, refreshKey string) (Maker, error) {
	bytes1, err := base64.StdEncoding.DecodeString(accessKey)
	if err != nil {
		return nil, err
	}
	bytes2, err := base64.StdEncoding.DecodeString(refreshKey)
	if err != nil {
		return nil, err
	}
	if len(bytes1) != chacha20poly1305.KeySize && len(bytes2) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid symmetric key size %d: must be %d bytes", len(accessKey), chacha20poly1305.KeySize)
	}
	maker := &PasetoMaker{
		paseto:     paseto.NewV2(),
		AccessKey:  []byte(bytes1),
		RefreshKey: []byte(bytes2),
	}
	return maker, nil
}

func (maker *PasetoMaker) CreateAccessToken(email string, isAdmin bool, duration time.Duration) (*Payload, error) {
	payload, err := NewPayload(email, isAdmin, duration)
	if err != nil {
		return nil, err
	}

	token, err := maker.paseto.Encrypt(maker.AccessKey, payload, nil)
	if err != nil {
		return nil, err
	}

	payload.Token = token
	return payload, nil
}
func (maker *PasetoMaker) VerifyAccessToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := maker.paseto.Decrypt(token, maker.AccessKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if err := payload.Valid(); err != nil {
		return nil, err
	}
	payload.Token = token

	return payload, nil
}

func (maker *PasetoMaker) CreateRefreshToken(email string, isAdmin bool, duration time.Duration) (*Payload, error) {
	payload, err := NewPayload(email, isAdmin, duration)
	if err != nil {
		return nil, err
	}

	token, err := maker.paseto.Encrypt(maker.RefreshKey, payload, nil)
	if err != nil {
		return nil, err
	}

	payload.Token = token
	return payload, nil
}
func (maker *PasetoMaker) VerifyRefreshToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := maker.paseto.Decrypt(token, maker.RefreshKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if err := payload.Valid(); err != nil {
		return nil, err
	}
	payload.Token = token

	return payload, nil
}
