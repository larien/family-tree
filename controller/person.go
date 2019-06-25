package controller

import (
	"log"
	r "github.com/larien/family-tree/repository"
	"github.com/larien/family-tree/entity"
)

// newPersonController applies person's Repository layer in Controller so
// that this layer can make use of Contract methods from Controller.
func newPersonController() *Person {
	log.Println("Person controller started")
	return &Person{}
}

// Person defines the object that contains methods from Repository layer.
type Person struct {
	Repository r.PersonRepository
}

// PersonController defines the method available from Person Controller
// domain to be used by external layers.
type PersonController interface {
	FindAll() ([]*entity.Person, error)
}

// FindAll returns all registered People.
func (p *Person) FindAll() ([]*entity.Person, error){
	return p.Repository.RetrieveAll(), nil
}