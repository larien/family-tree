package controller

import (
	"log"
	r "github.com/larien/family-tree/repository"
	"github.com/larien/family-tree/entity"
)

// newPersonController applies person's Repository layer in Controller so
// that this layer can make use of Contract methods from Controller.
func newPersonController(person r.PersonRepository) *Person {
	log.Println("Person controller started")
	return &Person{Repository: person}
}

// Person defines the object that contains methods from Repository layer.
type Person struct {
	Repository r.PersonRepository
}

// PersonController defines the method available from Person Controller
// domain to be used by external layers.
type PersonController interface {
	FindAll() ([]entity.Person, error)
	Find(string) (*entity.Person, error)
	Add([]entity.Person) error
}

// FindAll returns all registered People.
func (p *Person) FindAll() ([]entity.Person, error){
	return p.Repository.RetrieveAll()
}

// Find returns the Person data registered.
func (p *Person) Find(name string) (*entity.Person, error){
	return p.Repository.Retrieve(name)
}

// Add requests People and their relationships to be registered in the database.
func (p *Person) Add(people []entity.Person) error {
	for _, person := range people {
		log.Printf("Registering %s", person.Name)
		retrievedPerson, err := p.Repository.Retrieve(person.Name)
		if err != nil {return err}
		
		if retrievedPerson == nil {
			if err := p.Repository.Add(person.Name); err != nil {
				return err
			}
		}
		log.Printf("Registering %s's parents", person.Name)
		for _, parent := range person.Parents {
			if relationshipExists(parent, person.Parents){
				continue
			}
			retrievedParent, err := p.Repository.Retrieve(parent)
			if err != nil {return err }

			if retrievedParent == nil {
				if err := p.Repository.Add(parent); err != nil {
					return err
				}
			}

			err = p.Repository.Parent(parent, person.Name)
			if err != nil {return err }
		}
		
		log.Printf("Registering %s's children", person.Name)
		for _, child := range person.Children {
			if relationshipExists(child, person.Children){
				continue
			}
			retrievedChild, err := p.Repository.Retrieve(child)
			if err != nil {return err }

			if retrievedChild == nil {
				if err := p.Repository.Add(child); err != nil {
					return err
				}
			}

			err = p.Repository.Parent(person.Name, child)
			if err != nil {return err }
		}
		log.Printf("Registered %s", person.Name)
	}

	return nil
}

// relatinshopExists verify if the relationship already exists
// in the Person's data to prevent them to be duplicated.
func relationshipExists(newName string, names []string) bool {
	for _, name := range names {
		if newName == name {
			return false
		}
	}
	return true
}