package utility

import (
	"GoAddressBook/models"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
)

func RequestBodyValidator(contact models.Contact) error {
	if len(contact.FirstName) == 0 || len(contact.LastName) == 0 || len(contact.PhoneNumber) == 0 {
		return errors.New("given inputs are invalid")
	}
	return nil
}

func ParseValidatorErrMessage(err error) error {
	description := "Invalid request"
	if fieldErrors, ok := err.(validator.ValidationErrors); ok {
		err := fieldErrors[0]
		description = fmt.Sprintf("Invalid %s provided: %s", err.Field(), err.Value())
	}
	return errors.New(description)
}
