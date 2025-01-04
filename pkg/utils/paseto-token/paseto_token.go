package paseto_token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	Paseto       *paseto.V2
	SymmetricKey []byte
}

var TokenMaker *PasetoMaker

func NewPaseto(symmetricKey string) (*PasetoMaker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("SymmetricKey too short should be: %v", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		Paseto:       paseto.NewV2(),
		SymmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

func (maker *PasetoMaker) CreateToken(userId string, duration time.Duration) (string, error) {
	payload, err := NewPayload(userId, duration)
	if err != nil {
		return "", err
	}

	return maker.Paseto.Encrypt(maker.SymmetricKey, payload, nil)
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.Paseto.Decrypt(token, maker.SymmetricKey, payload, nil)
	if err != nil {
		return nil, err
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
