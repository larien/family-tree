package repository

import (
	"log"
)

// newPersonRepository applies database connection in Person
// Repository domain.
func newPersonRepository() *Person {
	log.Println("Person repository started")
	return &Person{}
}

// Person contains the database connection. Methods from Person Repository
// domain must implement this object in order to have access to database.
type Person struct {}

// PersonRepository defines the method available from Person Repository
// domain to be used by external layers.
type PersonRepository interface {}