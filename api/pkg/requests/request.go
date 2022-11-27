package requests

import (
	"fmt"
	"net/url"
	"strings"
	"unicode"

	"github.com/nyaruka/phonenumbers"
)

type request struct{}

func (request *request) sanitizeString(value string) string {
	return strings.TrimSpace(value)
}

func (request *request) sanitizePhoneNumber(value string) string {
	value = strings.TrimRight(value, " ")
	if len(value) > 0 && value[0] == ' ' {
		value = strings.Replace(value, " ", "+", 1)
	}

	if !strings.HasPrefix(value, "+") && request.isDigits(value) && len(value) > 9 {
		value = "+" + value
	}

	if number, err := phonenumbers.Parse(value, phonenumbers.UNKNOWN_REGION); err == nil {
		value = phonenumbers.Format(number, phonenumbers.E164)
	}

	return value
}

func (request *request) baseURL(value string) string {
	u, _ := url.Parse(value)
	return fmt.Sprintf("%s://%s", u.Scheme, u.Host)
}

func (request *request) isDigits(value string) bool {
	for _, c := range value {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}
