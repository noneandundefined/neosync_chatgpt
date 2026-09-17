package middleware

import (
	"context"
	"neomatica/neosync/infra/locale"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func LanguageMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lang := r.Header.Get("Accept-Language")
			if lang == "" {
				lang = "en"
			} else {
				lang = strings.Split(lang, ",")[0]
				lang = strings.Split(lang, "-")[0]
			}

			tr := locale.NewTranslator(lang)

			//nolint
			ctx := context.WithValue(r.Context(), "translator", tr)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func TranslatorFromContext(ctx context.Context) locale.Translator {
	if v := ctx.Value("translator"); v != nil {
		if tr, ok := v.(locale.Translator); ok {
			return tr
		}
	}

	return locale.NewTranslator("en")
}
