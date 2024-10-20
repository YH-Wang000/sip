package sipmsg

import (
	"log/slog"

	"sip/pkg/log"

	"github.com/go-playground/validator/v10"
)

func init() {
	err := Validator.RegisterValidation("checkSipVersion", checkSipVersion)
	if err != nil {
		log.Warn("Failed to register validation for sip version", slog.Any("error", err))
	}
}

var Validator = validator.New(validator.WithRequiredStructEnabled())

func checkSipVersion(fl validator.FieldLevel) bool {
	return fl.Field().String() == DefaultSipVersion
}
