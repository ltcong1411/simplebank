package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

// Payload contains the payload data of a token
type Payload struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	IssuedAt time.Time `json:"issued_at"`
	ExpireAt time.Time `json:"expire_at"`
}

// NewPayload creates a new payload with the username and duration
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Payload{
		ID:       tokenID,
		Username: username,
		IssuedAt: time.Now(),
		ExpireAt: time.Now().Add(duration),
	}, nil
}

func (p Payload) Valid() error {
	if time.Now().After(p.ExpireAt) {
		return ErrExpiredToken
	}
	return nil
}
