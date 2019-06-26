package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
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
	FamilyTree(string) ([]entity.FamilyTree, error)
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

// FamilyTree returns the Person's family tree.
func (p *Person) FamilyTree(name string) ([]entity.FamilyTree, error){
	log.Println("Getting family tree")

	people, err := p.Repository.RetrieveAll()
	if err != nil {return []entity.FamilyTree{}, err}

	err = dump(people)
	if err != nil {return []entity.FamilyTree{}, err}

	// 2 - get people without children
		// 2.1 - verify if name is between these people
		// 2.2 - if so, parse all relationships to FamilyTree
		// 2.3 if not, remove all people without children
			// repeat 2
	
	err = p.Person.Clear()
	if err != nil {return []entity.FamilyTree{}, err}

	// restore database
	err = removeDump()
	if err != nil {return []entity.FamilyTree{}, err}

	return []entity.FamilyTree{}, nil
}

func dump(people []entity.Person) error {
	filename := "dump.json"

	peopleJSON, err := json.Marshal(people)
	if err != nil {return err}

	err = ioutil.WriteFile(filename, peopleJSON, 0644)
	if err != nil {return err}

	log.Printf("Dump saved to %s", filename)

	return nil
}

func removeDump() error {
	return os.Remove("dump.json")
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