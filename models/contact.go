package models

// Contacts represents a contact and all its data in the address book
type Contact struct {
	FirstName    string `json:"first_name" validate:"omitempty,firstNameFormat"`
	LastName     string `json:"last_name" validate:"omitempty,lastNameFormat"`
	EmailAddress string `json:"email_address" validate:"omitempty,emailFormat"`
	PhoneNumber  string `json:"phone_number" validate:"omitempty,phoneNumberFormat"`
	Addresses    string `json:"address,omitempty"`
}

// Address represents a physical address for a contact
type Address struct {
	// The type of the address : personal, professional...
	Type string `json:"type"`

	// Next fields are precisions for the address
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
	Country string `json:"country"`
}
