package paseto_token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserId    string    `json:"userId"`
	CreatedAt time.Time `json:"created_at"`
	ExpiryAt  time.Time `json:"expiry_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenId,
		UserId:    username,
		CreatedAt: time.Now(),
		ExpiryAt:  time.Now().Add(duration * time.Minute),
	}

	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiryAt) {
		return errors.New("Token has expired")
	}
	return nil
}
