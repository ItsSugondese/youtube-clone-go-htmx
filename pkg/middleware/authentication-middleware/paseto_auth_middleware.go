package authentication_middleware

import (
	response_status_enum "youtube-clone/enums/interface-enums/response/response-status-enum"
	global_gin_context "youtube-clone/global/global-gin-context"
	globaldto "youtube-clone/global/global_dto"
	paseto_token "youtube-clone/pkg/utils/paseto-token"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeaderKey        = "Authorization"
	authorizationHeaderBearerType = "bearer"
)

func PasetoAuthMiddleware(maker paseto_token.PasetoMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		global_gin_context.GlobalGinContext.Context = ctx
		contextExtractToken, contextExtractTokenError := ExtractPasetoTokenFromHeaderContext(ctx)

		if contextExtractTokenError != nil {
			panic(contextExtractTokenError)
		}
		//token, extractTokenError := ExtractPasetoTokenFromHeader()
		//
		//if extractTokenError != nil {
		//	panic(extractTokenError)
		//}
		_, err := maker.VerifyToken(contextExtractToken)
		if err != nil {
			response := &globaldto.ApiResponse{
				Status:  response_status_enum.Fail(),
				Message: "Access Token Not Valid",
				Data:    []string{"Access Token Not Valid"},
			}
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		//context.WithValue(ctx.Request.Context(), "userToken", token)
		ctx.Next()
	}
}

func ExtractPasetoTokenFromHeader() (string, error) {
	if global_gin_context.GlobalGinContext.Context == nil {
		return "", errors.New("GinContext Context is nil")
	}
	authHeader := global_gin_context.GlobalGinContext.Context.GetHeader(authorizationHeaderKey)
	if authHeader == "" {
		return "", errors.New("No header was passed")
	}

	fields := strings.Fields(authHeader)
	if len(fields) != 2 {
		return "", errors.New("Invalid or Missing Bearer Token")
	}

	authType := fields[0]
	if strings.ToLower(authType) != authorizationHeaderBearerType {
		return "", errors.New("Authorization Type Not Supported")
	}

	return fields[1], nil
}

func ExtractPasetoTokenFromHeaderContext(ctx *gin.Context) (string, error) {
	authHeader := ctx.GetHeader(authorizationHeaderKey)
	if authHeader == "" {
		return "", errors.New("No header was passed")
	}

	fields := strings.Fields(authHeader)
	if len(fields) != 2 {
		return "", errors.New("Invalid or Missing Bearer Token")
	}

	authType := fields[0]
	if strings.ToLower(authType) != authorizationHeaderBearerType {
		return "", errors.New("Authorization Type Not Supported")
	}

	return fields[1], nil
}
