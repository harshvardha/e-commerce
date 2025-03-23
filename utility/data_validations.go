package utility

import "regexp"

func Validate_Phonenumber(phoneNumber string) bool {
	phoneNumberRegex := regexp.MustCompile(`^(?:(?:\+91|0)?[ -]?)?(?:(?:\d{2,4}[ -]?\d{6,8})|(?:\d{10}))$`)
	return phoneNumberRegex.MatchString(phoneNumber)
}
