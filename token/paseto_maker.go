package token

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto        *paseto.V2
	symmetrickKey []byte
}

func NewPasetoMaker(symmetrickKey string) (Maker, error) {
	if len(symmetrickKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size")
	}
	maker := &PasetoMaker{
		paseto:        paseto.NewV2(),
		symmetrickKey: []byte(symmetrickKey),
	}
	return maker, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	return maker.paseto.Encrypt(maker.symmetrickKey, payload, nil)
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := maker.paseto.Decrypt(token, maker.symmetrickKey, payload, nil)
	if err != nil {
		return nil, ErrorInvalidToken
	}
	err = payload.Valid()
	if err != nil {
		return nil, err
	}
	return payload, nil
}
