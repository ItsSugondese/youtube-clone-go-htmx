package localization

import (
	"fmt"
	"log"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/text/language"
)

type LocalizationManager struct {
	Localizer *i18n.Localizer
}

var GlobalLocalizationManager *LocalizationManager

func InitBundle() *i18n.Bundle {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// Load the translation files
	if _, err := bundle.LoadMessageFile("locales/en.toml"); err != nil {
		log.Fatalf("Failed to load locales/en.toml: %v", err)
	}
	if _, err := bundle.LoadMessageFile("locales/ne-NP.toml"); err != nil {
		log.Fatalf("Failed to load locales/ne-NP.toml: %v", err)
	}
	return bundle
}

func GetLocalizedMessage(messageID string, templateData map[string]interface{}) string {
	//localizer, exists := c.Get("localizer")
	//if !exists {
	//	panic("Localization error")
	//}

	localizer := GlobalLocalizationManager.Localizer

	if localizer == nil {
		panic("Localization error")
	}

	// Localize the message
	localizedMessage, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
		DefaultMessage: &i18n.Message{
			ID:    messageID,
			Other: "Message Not found",
		},
	})

	if err != nil {
		panic(fmt.Sprintf("Error localizing message: %v", err))
	}

	return localizedMessage
}

func InitLocalizationManager() {
	GlobalLocalizationManager = &LocalizationManager{}
}
