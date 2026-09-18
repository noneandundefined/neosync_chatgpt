package locale

import (
	"neomatica/neosync/config"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

type Translator interface {
	T(key string) string
	TErr(key string) string
	GetLang() string
}

func NewTranslator(lang string) Translator {
	return &translator{
		lang: lang,
		loc:  i18n.NewLocalizer(bundle, lang),
	}
}

type translator struct {
	lang string
	loc  *i18n.Localizer
}

func (t *translator) T(key string) string {
	msg, err := t.loc.Localize(&i18n.LocalizeConfig{MessageID: key})
	if err != nil {
		return key
	}

	return msg
}

func (t *translator) TErr(key string) string {
	return strings.ToLower(t.T(key))
}

func (t *translator) GetLang() string {
	return t.lang
}

func InitI18n() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", config.JSON.Unmarshal)

	bundle.MustLoadMessageFile("infra/locale/languages/en.json")
	bundle.MustLoadMessageFile("infra/locale/languages/es.json")
	bundle.MustLoadMessageFile("infra/locale/languages/ru.json")
}
