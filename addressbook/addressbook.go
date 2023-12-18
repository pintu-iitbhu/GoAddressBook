package addressbook

import (
	"GoAddressBook/constants"
	"GoAddressBook/models"
	"GoAddressBook/utility"
	"encoding/json"
	"fmt"
	"github.com/sagikazarmark/slog-shim"
	"os"
	"strings"
	"sync"
)

type AddressBook struct {
	Contacts   map[string]models.Contact `json:"contacts,omitempty"`    // Using a map for quick lookups
	NameIndex  map[string][]string       `json:"name_index,omitempty"`  // Index for name search
	PhoneIndex map[string]string         `json:"phone_index,omitempty"` // Index for phone search
	mutex      sync.RWMutex              // Mutex for concurrent access
}

// NewAddressBook creates a new AddressBook instance
func NewAddressBook() *AddressBook {
	return &AddressBook{
		Contacts:   make(map[string]models.Contact),
		NameIndex:  make(map[string][]string),
		PhoneIndex: make(map[string]string),
		mutex:      sync.RWMutex{},
	}
}

// LoadFromFile loads the address book from the JSON file
func (ab *AddressBook) LoadFromFile() error {
	data, err := os.ReadFile(constants.AddressBookFilePath)
	if err != nil {
		slog.Info("failed to read address book json file", err)
		return err
	}
	if err = json.Unmarshal(data, ab); err != nil {
		slog.Info("failed to unmarshal contacts details", err)
		return err
	}
	return nil
}

// AddContact CreateContact add a new contact into the book
func (ab *AddressBook) AddContact(contact models.Contact) {
	ab.mutex.Lock()
	defer ab.mutex.Unlock()

	key := ab.generateKey(contact.FirstName, contact.LastName, contact.PhoneNumber)
	ab.Contacts[key] = contact

	// Update name index
	nameKey := ab.generateNameKey(contact.FirstName, contact.LastName)
	ab.NameIndex[nameKey] = append(ab.NameIndex[nameKey], key)

	// Update phone index
	ab.PhoneIndex[contact.PhoneNumber] = key

	// Save to file
	ab.saveToFile()
}

// saveToFile saves the address book to the JSON file
func (ab *AddressBook) saveToFile() {
	err := os.Truncate(constants.AddressBookFilePath, 0)
	if err != nil {
		slog.Info("Failed to truncate address book file", err)
		return
	}

	file, _ := os.OpenFile(constants.AddressBookFilePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	data, err := json.MarshalIndent(ab, "", "    ")
	if err != nil {
		slog.Info("failed to marshal contacts details", err)
		return
	}

	_, err = file.Write(data)
	if err != nil {
		slog.Info("Error writing to file:", err)
		return
	}
	err = file.Close()
	return
}

// ListAllContacts returns a list of all the contacts in a pretty way
func (ab *AddressBook) ListAllContacts() []string {
	ab.mutex.RLock()
	defer ab.mutex.RUnlock()

	var contacts []string
	for _, actualContact := range ab.Contacts {
		actualContactByte, _ := json.Marshal(actualContact)
		contacts = append(contacts, string(actualContactByte))
	}
	return contacts
}

// SearchByName searches for a contact by name
func (ab *AddressBook) SearchByName(name string) []models.Contact {
	ab.mutex.RLock()
	defer ab.mutex.RUnlock()

	firstName, _, lastName := utility.GetFirstMiddleAndLastNamesFromFullName(name)
	nameKey := ab.generateNameKey(firstName, lastName)
	keys, found := ab.NameIndex[nameKey]
	if !found {
		return nil
	}

	var results []models.Contact
	for _, key := range keys {
		contact, _ := ab.Contacts[key]
		results = append(results, contact)
	}

	return results
}

// SearchByPhoneNumber searches for a contact by phone number
func (ab *AddressBook) SearchByPhoneNumber(phoneNumber string) (models.Contact, bool) {
	ab.mutex.RLock()
	defer ab.mutex.RUnlock()

	key, found := ab.PhoneIndex[phoneNumber]
	if !found {
		return models.Contact{}, false
	}

	contact, found := ab.Contacts[key]
	return contact, found
}

// generateKey generates a unique key for a contact based on first name, last name, and phone number
func (ab *AddressBook) generateKey(firstName, lastName, phoneNumber string) string {

	return fmt.Sprintf("%s-%s-%s", firstName, lastName, phoneNumber)
}

// generateNameKey generates a unique key for indexing names
func (ab *AddressBook) generateNameKey(firstName, lastName string) string {
	return fmt.Sprintf("%s-%s", strings.ToLower(firstName), strings.ToLower(lastName))
}
