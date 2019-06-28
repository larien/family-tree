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
	Add([]entity.Person) error
	Ascendancy(string) ([]entity.Person, error)
	Find(string) (*entity.Person, error)
	FindAll() ([]entity.Person, error)
	Clear() error
}

// Add requests People and their relationships to be registered in the database.
func (p *Person) Add(people []entity.Person) error {
	for _, person := range people {
		err := p.RegisterPerson(person.Name)
		if err != nil {return err}

		err = p.RegisterParents(person.Name, person.Parents)
		if err != nil {return err}

		err = p.RegisterChildren(person.Name, person.Children)
		if err != nil {return err}

		log.Printf("Registered %s", person.Name)
	}

	return nil
}

// Ascendancy returns the Person's family tree. This algorithm works as
// explained below:
// We check if the Person where the ascendancy begins from exists.
// Every People in the tree is retrieved in order to create a backup
// file containing the current data. This is made because the data
// inside the database will be changed.
// We have to find a way to navigate between the parentship levels
// in order to get ascendancy, so every children with no children
// and with parents is deleted till the Person has no children.
// This way, the generated graph is the connection between the
// person and its ascendants.
// Therefore, we do a search to get their ascendants and restore data.
func (p *Person) Ascendancy(name string) ([]entity.Person, error){
	log.Printf("Getting %s's ascendancy", name)

	filename := "dump.json"

	person, err := p.Repository.Retrieve(name)
	if err != nil {return []entity.Person{}, err}

	if person == nil {
		return []entity.Person{}, fmt.Errorf("%s wasn't found", name)
	}

	err = p.Repository.Backup(filename)
	if err != nil {return []entity.Person{}, err}

	err = p.Ascend(name)
	if err != nil {return []entity.Person{}, err}

	connectedNames, err := p.Repository.Connected(name)
	if err != nil {return []entity.Person{}, err} 

	err = p.Restore(filename)
	if err != nil {return []entity.Person{}, err}
	
	return p.Ascendants(connectedNames)
}

// Find returns the Person data registered.
func (p *Person) Find(name string) (*entity.Person, error){
	log.Printf("Finding %s\n", name)

	return p.Repository.Retrieve(name)
}

// FindAll returns all registered People.
func (p *Person) FindAll() ([]entity.Person, error){
	log.Println("Finding all People")

	return p.Repository.RetrieveAll()
}

// Clear is a helper function that cleans the database.
func (p *Person) Clear() error {
	log.Println("Cleaning database")

	return p.Repository.Clear()
}

// Register register the relationship to the Person.
func (p *Person) Register(related string, person *entity.Person) error {
	retrieved, err := p.Repository.Retrieve(related)
	if err != nil {return err }

	if retrieved == nil {
		if err := p.Repository.Add(related); err != nil {
			return err
		}
	}
	return nil
}

// RegisterPerson registers the Person if it doesn't exist in
// the database.
func (p *Person) RegisterPerson(name string) error {
	log.Printf("Registering %s", name)
	retrievedPerson, err := p.Repository.Retrieve(name)
	if err != nil {return err}
		
	if retrievedPerson == nil {
		if err := p.Repository.Add(name); err != nil {
			return err
		}
	}
	return nil
}

// RegisterParents register the Person's parents.
func (p *Person) RegisterParents(name string, parents []string) error {
	log.Printf("Registering %s's parents", name)
		for _, parent := range parents {
			retrievedPerson, err := p.Repository.Retrieve(name)
			if err != nil {return err}
		
			if relationshipExists(parent, retrievedPerson.Parents){
				continue
			}

			err = p.Register(parent, retrievedPerson)
			if err != nil {return err}

			err = p.Repository.Parent(parent, name)
			if err != nil {return err}
		}
	return nil
}

// RegisterChildren register the Person's children.
func (p *Person) RegisterChildren(name string, children []string) error {
	log.Printf("Registering %s's children", name)
	for _, child := range children {
		retrievedPerson, err := p.Repository.Retrieve(name)
		if err != nil {return err}

		if relationshipExists(child, retrievedPerson.Children){
			continue
		}

		err = p.Register(child, retrievedPerson)
		if err != nil {return err}

		err = p.Repository.Parent(name, child)
			if err != nil {return err}
	}
	return nil
}

// Ascend removes the lowest-level relationships in order to
// ascent the parentship tree till it gets to the Person requested.
func (p *Person) Ascend(name string) error {
	for {
		children, err := p.Repository.Children(name)
		if err != nil {return err}
		
		if children == nil {
			break
		}

		err = p.Repository.DeleteWithoutChildren()
		if err != nil {return err}
	}

	return nil
}

// Ascendants gets the People's connected relationships.
func (p *Person) Ascendants(connectedNames []string) (ascendants []entity.Person, err error) {
	for _, connectedName := range connectedNames {
		person, err := p.Repository.Retrieve(connectedName)
		if err != nil {return []entity.Person{}, err}

		ascendants = append(ascendants, *person)
	}
	return
}

// Restore restores People from the dump file.
func (p *Person) Restore(filename string) error {
	err := p.Repository.Clear()
	if err != nil {return err}

	people, err := readDump(filename)
	if err != nil {return err}

	err = p.Add(people)
	if err != nil {return err}

	log.Printf("Database restored from %s", filename)

	err = os.Remove(filename)
	if err != nil {return err}

	return nil
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

// relationshipExists verify if the relationship already exists
// in the Person's data to prevent them to be duplicated.
func relationshipExists(newName string, names []string) bool {
	for _, name := range names {
		if newName == name {
			return true
		}
	}
	return false
}