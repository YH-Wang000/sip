package sipmsg

import (
	"log/slog"
	"slices"

	"sip/pkg/log"

	"github.com/go-playground/validator/v10"
)

func init() {
	err := Validator.RegisterValidation("checkSipVersion", checkSipVersion)
	if err != nil {
		log.Warn("Failed to register validation for sip version", slog.Any("error", err))
	}
	err = Validator.RegisterValidation("checkSipRequestMethod", checkSipRequestMethod)
	if err != nil {
		log.Warn("Failed to register validation for sip request method enum", slog.Any("error", err))
	}
	err = Validator.RegisterValidation("checkSipScheme", checkSipScheme)
	if err != nil {
		log.Warn("Failed to register validation for sip scheme enum", slog.Any("error", err))
	}
}

var Validator = validator.New(validator.WithRequiredStructEnabled())

func checkSipVersion(fl validator.FieldLevel) bool {
	return fl.Field().String() == DefaultSipVersion
}

func checkSipRequestMethod(fl validator.FieldLevel) bool {
	return SipRequestMethodSet[MethodEnum(fl.Field().String())]
}

func checkSipScheme(fl validator.FieldLevel) bool {
	return slices.Contains([]SipSchemeEnum{Sip, Sips, Tel}, SipSchemeEnum(fl.Field().String()))
}
