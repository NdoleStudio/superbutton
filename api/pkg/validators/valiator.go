package validators

import (
	"fmt"

	"github.com/nyaruka/phonenumbers"
	"github.com/thedevsaddam/govalidator"
)

const (
	phoneNumberRule = "phoneNumber"
)

func init() {
	// custom rules to take fixed length word.
	// e.g: max_word:5 will throw error if the field contains more than 5 words
	govalidator.AddCustomRule(phoneNumberRule, func(field string, rule string, message string, value interface{}) error {
		phoneNumber, ok := value.(string)
		if !ok {
			return fmt.Errorf("the %s field must be a valid E.164 phone number e.g +18005550199", field)
		}

		_, err := phonenumbers.Parse(phoneNumber, phonenumbers.UNKNOWN_REGION)
		if err != nil {
			return fmt.Errorf("the %s field must be a valid E.164 phone number e.g +18005550199", field)
		}

		return nil
	})
}
