package dto

type OAuth2ClientRequest struct {
	ClientID     string `json:"clientId" binding:"required"`
	ClientSecret string `json:"clientSecret" binding:"required"`
	Domain       string `json:"domain" binding:"required"`
}
