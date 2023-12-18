package utility

import (
	"GoAddressBook/models"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

func RequestBodyValidator(contact models.Contact) error {
	if len(contact.FirstName) == 0 || len(contact.LastName) == 0 || len(contact.PhoneNumber) == 0 || len(contact.EmailAddress) == 0 {
		return errors.New("given inputs are invalid")
	}
	return nil
}

func ParseValidatorErrMessage(err error) error {
	description := "Invalid request"
	var fieldErrors validator.ValidationErrors
	if errors.As(err, &fieldErrors) {
		err := fieldErrors[0]
		description = fmt.Sprintf("Invalid %s provided: %s", err.Field(), err.Value())
	}
	return errors.New(description)
}

func GetFirstMiddleAndLastNamesFromFullName(fullName string) (string, string, string) {
	fullName = removeMultiSpaceFromFullName(fullName)
	splitNames := strings.SplitN(fullName, " ", 3)

	var firstName string
	var middleName string
	var lastName string

	if len(splitNames) >= 1 {
		firstName = splitNames[0]
	}

	if len(splitNames) == 3 {
		middleName = splitNames[1]
	}

	if len(splitNames) >= 2 {
		if len(splitNames) >= 3 {
			lastName = splitNames[2]
		} else {
			lastName = splitNames[1]
		}
	}
	return firstName, middleName, lastName
}

func removeMultiSpaceFromFullName(s string) string {
	str := strings.Fields(s)
	var out string
	for i := 0; i < len(str); i++ {
		out += str[i] + " "
	}
	return strings.TrimSpace(out)
}
