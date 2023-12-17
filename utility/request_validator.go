package utility

import (
	"GoAddressBook/constants"
	"github.com/go-playground/validator/v10"
	"github.com/sagikazarmark/slog-shim"
	"reflect"
	"regexp"
	"strings"
)

func NewValidator() *validator.Validate {
	requestValidator := validator.New()
	requestValidator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	var err error
	defer func() {
		if err != nil {
			panic("Validator registration failed " + err.Error())
		}
	}()

	err = requestValidator.RegisterValidation("firstNameFormat", FullNameFormatValidator)
	if err != nil {
		slog.Info("Error while registering custom validator func FullNameFormatValidator %s\n", err.Error())
		return nil
	}

	err = requestValidator.RegisterValidation("lastNameFormat", FullNameFormatValidator)
	if err != nil {
		slog.Info("Error while registering custom validator func FullNameFormatValidator %s\n", err.Error())
		return nil
	}
	err = requestValidator.RegisterValidation("emailFormat", EmailFormatValidator)
	if err != nil {
		slog.Info("Error while registering custom validator func EmailFormatValidator %s\n", err.Error())
		return nil
	}
	err = requestValidator.RegisterValidation("phoneNumberFormat", PhoneNumberFormatValidator)
	if err != nil {
		slog.Info("Error while registering custom validator func PhoneNumberFormatValidator %s\n", err.Error())
		return nil
	}

	return requestValidator
}

func FullNameFormatValidator(fl validator.FieldLevel) bool {
	fullName := fl.Field().String()
	reg, err := regexp.Compile(constants.SalutationRegex)
	if err != nil {
		return false
	}

	processedString := reg.ReplaceAllString(fullName, "")
	if len(processedString) > 750 {
		slog.Info("fullName Length should be from 1 to 750")
		return false
	}

	/*	TODO : might introduce in future. Removing currently to cater to requests coming with . in prod same as old PD service
		isStringAlphabetic := regexp.MustCompile(constants.AlphabetOnlyRegex).MatchString
	*/
	// remove wide space from full name
	str := strings.Replace(processedString, " ", "", -1)
	/*
		if !(isStringAlphabetic(str)) {
			return false
		}
	*/
	if len(strings.Split(processedString, "")) <= 0 {
		slog.Info("FullName length is less than 0 or not contains word, after processing : %s\n", str)
		return false
	}
	return true
}

func EmailFormatValidator(fl validator.FieldLevel) bool {
	fieldValue := fl.Field().String()
	return IsValidEmailAddress(fieldValue)
}
func IsValidEmailAddress(emailAddress string) bool {
	reg := regexp.MustCompile(constants.EmailRegex)
	tokenizedEmail := strings.Split(emailAddress, "@")
	isValidEmail := reg.MatchString(emailAddress)
	if !(isValidEmail) || len(tokenizedEmail[0]) <= 2 {
		return false
	}
	return true
}

func PhoneNumberFormatValidator(fl validator.FieldLevel) bool {
	fieldValue := fl.Field().String()
	reg := regexp.MustCompile(constants.PhoneNumberRegex)
	phoneNumber := reg.FindStringSubmatch(fieldValue)
	return phoneNumber != nil
}
