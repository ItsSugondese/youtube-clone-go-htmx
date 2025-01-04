package service

import (
	"os"
	"youtube-clone/enums/struct-enums/user_type_enums"
	"youtube-clone/internal/auth/dto"
	"github.com/go-oauth2/oauth2/v4/models"
	oauth2_setup "youtube-clone/config/oauth2-setup"
	"youtube-clone/global/global_var"
	user_navigator "youtube-clone/internal/user/user-navigator"
	"youtube-clone/pkg/utils"
	paseto_token "youtube-clone/pkg/utils/paseto-token"
	"github.com/gin-gonic/gin"
	"youtube-clone/pkg/common/database"
	"time"
)

func LoginService(authRequest dto.AuthRequest) dto.AuthResponse {
	switch authRequest.UserType {
	case user_type_enums.UserType.CUSTOMER, user_type_enums.UserType.ADMIN:
		return loginUser(authRequest)
	default:
		panic("invalid user type")
	}
}

func RegisterOAuth2ClientService(ctx *gin.Context, authDto dto.OAuth2ClientRequest) dto.OAuth2ClientResponse {
	client := models.Client{
		ID:     authDto.ClientID,
		Secret: authDto.ClientSecret,
		Domain: authDto.Domain,
	}

	// Check if the client already exists
	info, _ := oauth2_setup.OAuthServerDetails.ClientStore.GetByID(ctx, authDto.ClientID)

	if info != nil {

		// Client exists, delete it first
		result := database.DB.Exec("UPDATE "+global_var.OAuthClientTable+" SET secret = ?, domain = ? WHERE id = ?", authDto.ClientSecret,
			authDto.Domain, authDto.ClientID)

		// Check for errors and print rows affected
		if result.Error != nil {
			panic(result.Error)
		}
	} else {
		err := oauth2_setup.OAuthServerDetails.ClientStore.Create(&client)
		if err != nil {
			panic(err)
		}
	}

	// Add the updated or new client

	return dto.OAuth2ClientResponse{}
}

func loginUser(authRequest dto.AuthRequest) dto.AuthResponse {
	userDetails := user_navigator.FindUserByEmailService(authRequest.Email)

	err := utils.VerifyPassword(userDetails.Password, authRequest.Password)
	if err != nil {
		panic(err)
	}

	maker := *paseto_token.TokenMaker

	token, err := maker.CreateToken(userDetails.ID.String(), 21600)
	if err != nil {
		panic(err)
	}

	return dto.AuthResponse{
		Token: token,
	}
}

func createToken(userID string) (string, error) {
	privateKey := os.Getenv("ACCESS_TOKEN_PRIVATE_KEY")
	expireTime := os.Getenv("ACCESS_TOKEN_EXPIRED_IN")
	expireDuration, err := time.ParseDuration(expireTime)
	if err != nil {
		return "", err
	}

	return utils.CreateToken(expireDuration, userID, privateKey)
}
