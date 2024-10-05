package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserName  string    `json:"user_name"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

var (
	ErrorInvalidToken = errors.New("token is invalid")
	ErrorExpiredToken = errors.New("token is expired")
)

// Create a new payload with the given username and duration.
func NewPayload(userName string, duration time.Duration) (*Payload, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		ID:        id,
		UserName:  userName,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrorExpiredToken
	}
	return nil
}
