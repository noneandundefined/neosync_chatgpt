package httpx

import (
	"fmt"
	"neomatica/neosync/infra/locale"
	"regexp"

	"github.com/go-playground/validator"
)

var Validate *validator.Validate
var passwordRegex = regexp.MustCompile(`^[A-Za-z0-9!"#$%&'()*+\-./:;<=>?@[\\\]^_{|}~]{0,8}$`)

func init() {
	Validate = validator.New()

	_ = Validate.RegisterValidation("optional_uuid", func(fl validator.FieldLevel) bool {
		uuid := fl.Field().String()
		if uuid == "" {
			return true
		}

		return validator.New().Var(uuid, "uuid") == nil
	})

	_ = Validate.RegisterValidation("imei", func(fl validator.FieldLevel) bool {
		imei := fl.Field().String()
		if len(imei) != 15 {
			return false
		}

		matched, _ := regexp.MatchString(`^[0-9]{15}$`, imei)
		return matched
	})

	_ = Validate.RegisterValidation("dpass", func(fl validator.FieldLevel) bool {
		pass := fl.Field().String()
		if pass == "" {
			return true
		}

		return passwordRegex.MatchString(pass)
	})
}

func ValidateMsg(tr locale.Translator, err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			switch fieldError.Tag() {
			case "required":
				return fmt.Sprintf(tr.TErr("field-required"), fieldError.Field())
			case "len":
				return fmt.Sprintf(tr.TErr("field-len"), fieldError.Field(), fieldError.Param())
			case "numeric":
				return fmt.Sprintf(tr.TErr("field-numeric"), fieldError.Field())
			case "uuid":
				return fmt.Sprintf(tr.TErr("field-uuid"), fieldError.Field())
			case "gt":
				return fmt.Sprintf(tr.TErr("field-gt"), fieldError.Field(), fieldError.Param())
			case "max":
				return fmt.Sprintf(tr.TErr("field-max"), fieldError.Field(), fieldError.Param())
			case "min":
				return fmt.Sprintf(tr.TErr("field-min"), fieldError.Field(), fieldError.Param())
			case "optional_uuid":
				return fmt.Sprintf(tr.TErr("field-optional-uuid"), fieldError.Field())
			case "imei":
				return tr.TErr("field-optional-imei")
			case "dpass":
				return tr.TErr("requirements-pass")
			default:
				return fmt.Sprintf(tr.TErr("field-validation-error"), fieldError.Field())
			}
		}
	}

	return tr.TErr("request-validation-error")
}
