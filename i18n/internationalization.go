package i18n

import (
	"GoAddressBook/constants"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Internationalization struct {
	Localizer *i18n.Localizer
}

func NewI18nInstance(locale string) (instance Internationalization, err error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc(constants.TomlFileFormat, toml.Unmarshal)

	if _, err := bundle.LoadMessageFile(constants.EnTomlFilePath); err != nil {
		return instance, err
	}
	if _, err := bundle.LoadMessageFile(constants.FrTomlFilePath); err != nil {
		return instance, err
	}

	localizer := i18n.NewLocalizer(bundle, locale)
	instance.Localizer = localizer

	return
}

func (instance Internationalization) T(key string, params map[string]interface{}) (message string, err error) {
	message, err = instance.Localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: params,
	})
	return
}
