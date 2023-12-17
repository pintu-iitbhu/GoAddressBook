package cli

import (
	"GoAddressBook/addressbook"
	"GoAddressBook/constants"
	"GoAddressBook/i18n"
	"GoAddressBook/models"
	"GoAddressBook/utility"
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/chzyer/readline"
	"github.com/go-playground/validator/v10"
	"github.com/sagikazarmark/slog-shim"
	"github.com/spf13/viper"
)

// Cli structure represents the command line interface
type Cli struct {
	Book      *addressbook.AddressBook
	Reader    *readline.Instance
	I18n      *i18n.Internationalization
	Validator *validator.Validate
}

// NewCliInstance NewInstance returns an instance of the Cli structure
func NewCliInstance(book *addressbook.AddressBook) (cli *Cli, err error) {
	reader, err := readline.New("> ")
	if err != nil {
		return &Cli{}, err
	}
	locale := viper.GetString(constants.Locale)
	i18nInstance, err := i18n.NewI18nInstance(locale)
	return &Cli{
		Book:      book,
		Reader:    reader,
		I18n:      &i18nInstance,
		Validator: utility.NewValidator(),
	}, nil
}

// Menu displays and loops over the menu in the command line interface
func (instance *Cli) Menu() {
	openingString, _ := instance.I18n.T(constants.Opening, nil)
	println(openingString)

	actionsString, _ := instance.I18n.T(constants.Actions, nil)
	listString, _ := instance.I18n.T(constants.List, nil)
	createContact, _ := instance.I18n.T(constants.Create, nil)
	searchByPhone, _ := instance.I18n.T(constants.SearchByPhoneNumber, nil)
	searchByName, _ := instance.I18n.T(constants.SearchByName, nil)
	closeString, _ := instance.I18n.T(constants.Close, nil)
	unknownChoiceString, _ := instance.I18n.T(constants.UnknownChoice, nil)

	var choice int
	prompt := &survey.Select{
		Message: actionsString,
		Options: []string{
			createContact,
			searchByName,
			searchByPhone,
			listString,
			closeString,
		},
	}

	quit := false
	for !quit {
		_ = survey.AskOne(prompt, &choice)

		switch choice {
		case 0:
			instance.Create()
		case 1:
			instance.GetContactDetailsByName()
		case 2:
			instance.GetContactDetailsByPhoneNumber()
		case 3:
			instance.List()
		case 4:
			quit = true
		default:
			println(unknownChoiceString)
		}
	}

	closingString, _ := instance.I18n.T(constants.Closing, nil)
	println(closingString)
	_ = instance.Reader.Close()
}
func (instance *Cli) List() {
	listingString, _ := instance.I18n.T(constants.ContactsListing, nil)
	println(listingString)
	for _, contact := range instance.Book.ListAllContacts() {
		println(contact)
		println("--------->")
	}
}

// Create prompts the user to add a contact using the command line interface
func (instance *Cli) Create() {
	addingString, _ := instance.I18n.T(constants.ContactAdding, nil)
	firstNameString, _ := instance.I18n.T(constants.FirstName, nil)
	lastNameString, _ := instance.I18n.T(constants.LastName, nil)
	phoneNumber, _ := instance.I18n.T(constants.Phone, nil)
	eMailAddress, _ := instance.I18n.T(constants.Email, nil)
	addressDetails, _ := instance.I18n.T(constants.Address, nil)

	println(addingString)
	contact := models.Contact{
		FirstName:    instance.readLine(firstNameString),
		LastName:     instance.readLine(lastNameString),
		PhoneNumber:  instance.readLine(phoneNumber),
		EmailAddress: instance.readLine(eMailAddress),
		Addresses:    instance.readLine(addressDetails),
	}

	err := utility.RequestBodyValidator(contact)
	if err != nil {
		println("Requested data is invalid", "err", err)
		println("------->")
		return
	}

	validationErr := instance.Validator.Struct(contact)
	if validationErr != nil {
		println("failed to validate request", "err: ", utility.ParseValidatorErrMessage(validationErr).Error())
		println("------->")
		return
	}
	instance.Book.AddContact(contact)
	addedString, _ := instance.I18n.T(constants.ContactAdded, map[string]interface{}{
		constants.FirstName:   contact.FirstName,
		constants.LastName:    contact.LastName,
		constants.PhoneNumber: contact.PhoneNumber,
		constants.Email:       contact.EmailAddress,
		constants.Address:     constants.Address,
	})
	println(addedString)
}

func (instance *Cli) GetContactDetailsByName() {
	searchByName, _ := instance.I18n.T(constants.SearchByFullName, nil)
	name := instance.readLine(searchByName)
	contacts := instance.Book.SearchByName(name)
	if contacts == nil {
		println("Contacts not found for User: ", name)
		return
	}
	println("List of contact details of user by name: ")
	for _, actualContact := range contacts {
		actualContactByte, _ := json.Marshal(actualContact)
		actualContactByteString := string(actualContactByte)
		println(actualContactByteString)
		println("------------->")
	}
}

func (instance *Cli) GetContactDetailsByPhoneNumber() {
	searchByPhone, _ := instance.I18n.T(constants.SearchByPhoneNumber, nil)
	phone := instance.readLine(searchByPhone)
	actualContact, found := instance.Book.SearchByPhoneNumber(phone)
	if !found {
		fmt.Println("Contact details not found on Phone Number: ", phone)
		return
	}
	actualContactByte, _ := json.Marshal(actualContact)
	actualContactByteString := string(actualContactByte)
	println(actualContactByteString)
	println("------------->")
}

// Function to read a line and handle errors
func (instance *Cli) readLine(prompt string) string {
	println(prompt)
	line, err := instance.Reader.Readline()
	if err != nil {
		slog.Info("Error occurred while reading command line", "err:", err)
		panic(err)
	}
	return line
}
