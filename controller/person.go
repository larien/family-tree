package controller

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
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
	Ascendancy(string) ([]entity.Person, error)
	Add([]entity.Person) error
	Restore(string) (error)
}

// FindAll returns all registered People.
func (p *Person) FindAll() ([]entity.Person, error){
	return p.Repository.RetrieveAll()
}

// Find returns the Person data registered.
func (p *Person) Find(name string) (*entity.Person, error){
	return p.Repository.Retrieve(name)
}

// Ascendancy returns the Person's family tree.
func (p *Person) Ascendancy(name string) ([]entity.Person, error){
	log.Printf("Getting %s's ascendancy", name)

	filename := "dump.json"

	people, err := p.Repository.RetrieveAll()
	if err != nil {return []entity.Person{}, err}
	err = dump(people, filename)
	if err != nil {return []entity.Person{}, err}

	person, err := p.Repository.Retrieve(name)
	if err != nil {return []entity.Person{}, err}

	if person == nil {
		return []entity.Person{}, fmt.Errorf("%s wasn't found", name)
	}
	for {
		children, err := p.Repository.Children(name)
		if err != nil {return []entity.Person{}, err}
		if children == nil {
			break
		}

		err = p.Repository.DeleteWithoutChildren()
		if err != nil {return []entity.Person{}, err}
	}

	connectedNames, err := p.Repository.Connected(name)
	if err != nil {return []entity.Person{}, err} 

	err = p.Repository.Clear()
	if err != nil {return []entity.Person{}, err}
	err = p.Restore(filename)
	if err != nil {return []entity.Person{}, err}
	
	ascendants := []entity.Person{}
	for _, connectedName := range connectedNames {
		person, err := p.Repository.Retrieve(connectedName)
		if err != nil {return []entity.Person{}, err}

		ascendants = append(ascendants, *person)
	}
	
	err = removeDump()
	if err != nil {return []entity.Person{}, err}

	return ascendants, nil
}

// Restore restores People from the system from a dump file.
func (p *Person) Restore(filename string) error{
	people, err := readDump(filename)
	if err != nil {return err}

	return p.Add(people)
}

// readDump opens the dump file and restores it to the memory.
func readDump(filename string) ([]entity.Person, error) {
	jsonFile, err := os.Open(filename)
	if err != nil {return nil, err}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {return nil, err}

	var people []entity.Person

	err = json.Unmarshal(byteValue, &people)
	if err != nil {return nil, err}

	return people, nil
}

func dump(people []entity.Person, filename string) error {
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
			retrievedPerson, err := p.Repository.Retrieve(person.Name)
			if err != nil {return err}
		
			if relationshipExists(parent, retrievedPerson.Parents){
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
			retrievedPerson, err := p.Repository.Retrieve(person.Name)
			if err != nil {return err}

			if relationshipExists(child, retrievedPerson.Children){
				fmt.Println("Pais: ", retrievedPerson.Children)
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
			return true
		}
	}
	return false
}