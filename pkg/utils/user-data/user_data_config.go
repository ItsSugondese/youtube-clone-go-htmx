package user_data

import (
	authentication_middleware "youtube-clone/pkg/middleware/authentication-middleware"
	"youtube-clone/pkg/utils/paseto-token"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func DecryptToken(maker *paseto_token.PasetoMaker) (*paseto_token.Payload, error) {
	payload := &paseto_token.Payload{}

	token, extractTokenError := authentication_middleware.ExtractPasetoTokenFromHeader()

	if extractTokenError != nil {
		return payload, errors.New(extractTokenError.Error())
	}

	err := maker.Paseto.Decrypt(token, maker.SymmetricKey, payload, nil)
	if err != nil {
		return payload, errors.New(err.Error())

	}

	err = payload.Valid()
	if err != nil {
		return payload, errors.New(err.Error())
	}

	return payload, nil
}

func DecryptTokenContext(ctx *gin.Context, maker *paseto_token.PasetoMaker) (*paseto_token.Payload, error) {
	payload := &paseto_token.Payload{}

	token, extractTokenError := authentication_middleware.ExtractPasetoTokenFromHeaderContext(ctx)

	if extractTokenError != nil {
		return payload, errors.New(extractTokenError.Error())
	}

	err := maker.Paseto.Decrypt(token, maker.SymmetricKey, payload, nil)
	if err != nil {
		return payload, errors.New(err.Error())

	}

	err = payload.Valid()
	if err != nil {
		return payload, errors.New(err.Error())
	}

	return payload, nil
}

func GetUserIdContext(ctx *gin.Context) (string, error) {
	payload, err := DecryptTokenContext(ctx, paseto_token.TokenMaker)
	if err != nil {
		return "", err
	}
	return payload.UserId, nil
}

func GetUserIdErrorHandledContext(ctx *gin.Context) uuid.UUID {
	payload, err := DecryptTokenContext(ctx, paseto_token.TokenMaker)
	if err != nil {
		panic(err)
	}

	parsedUserId, parsingErr := uuid.Parse(payload.UserId)

	if parsingErr != nil {
		panic(parsingErr)
	}
	return parsedUserId
}
