package oauth2_setup

import (
	"context"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/jackc/pgx/v4"
	pg "github.com/vgarvardt/go-oauth2-pg/v4"
	"github.com/vgarvardt/go-pg-adapter/pgx4adapter"
	"os"
	"youtube-clone/global/global_var"
	baseUserModel "youtube-clone/internal/user/model"
	"youtube-clone/pkg/common/database"
	"youtube-clone/pkg/utils"
	"time"
)

var OAuthServerDetails *oAuthServerDetails

type oAuthServerDetails struct {
	Srv         *server.Server
	ClientStore *pg.ClientStore
}

func SetUpOAuth2() {

	pgxConn, _ := pgx.Connect(context.TODO(), os.Getenv("DB_URL"))

	manager := manage.NewDefaultManager()
	adapter := pgx4adapter.NewConn(pgxConn)

	tokenStore, _ := pg.NewTokenStore(adapter, pg.WithTokenStoreGCInterval(time.Minute))
	defer tokenStore.Close()

	clientStore, _ := pg.NewClientStore(adapter, pg.WithClientStoreTableName(global_var.OAuthClientTable))

	manager.MapTokenStorage(tokenStore)
	manager.MapClientStorage(clientStore)
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	manager.MustTokenStorage(store.NewMemoryTokenStore())

	srv := server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)
	manager.SetRefreshTokenCfg(manage.DefaultRefreshTokenCfg)

	srv.SetPasswordAuthorizationHandler(func(ctx context.Context, clientID, username, password string) (userID string, err error) {
		// Implement your own user authentication (e.g., verify credentials in your base_users table)
		var user baseUserModel.BaseUser
		err = database.DB.Where("email = ?", username).First(&user).Error
		if err != nil {
			return "", errors.ErrInvalidGrant // Return the error instead of panicking
		}

		err = utils.VerifyPassword(user.Password, password)
		if err != nil {
			return "", err // Return the error if password verification fails
		}

		return user.ID.String(), nil // Return the userID (email) and nil error if successful
	})

	//client := &models.Client{
	//	ID:     "client_id",
	//	Secret: "secret",
	//	Domain: "https://example.com",
	//}
	//ClientStore.Create(client)

	OAuthServerDetails = &oAuthServerDetails{
		Srv:         srv,
		ClientStore: clientStore,
	}
}
