package utility

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	Validate *validator.Validate
}

func ValidatePhonenumber(phoneNumber string) bool {
	phoneNumberRegex := regexp.MustCompile(`^(?:(?:\+91|0)?[ -]?)?(?:(?:\d{2,4}[ -]?\d{6,8})|(?:\d{10}))$`)
	return phoneNumberRegex.MatchString(phoneNumber)
}

func (validator *Validator) ValidateEmail(email string) bool {
	err := validator.Validate.Var(email, "required,email")
	return err == nil
}

func (validator *Validator) ValidatePassword(password string) bool {
	return validator.Validate.Var(password, "required,min=6,max=64,password") == nil
}
