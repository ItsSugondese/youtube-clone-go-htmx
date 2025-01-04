package lang_middleware

import (
	"youtube-clone/pkg/common/localization"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func LocalizationMiddleware(bundle *i18n.Bundle) gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.DefaultQuery("lang", language.English.String())
		localizer := i18n.NewLocalizer(bundle, lang)
		//c.Set("localizer", localizer)
		localization.GlobalLocalizationManager.Localizer = localizer
		c.Next()
	}
}
