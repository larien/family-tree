package repository

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"log"
	"fmt"
)

// newPersonRepository applies database connection in Person
// Repository domain.
func newPersonRepository(connection *DBConnection) *Person {
	log.Println("Person repository started")
	return &Person{DB: connection}
}

// Person contains the database connection. Methods from Person Repository
// domain must implement this object in order to have access to database.
type Person struct {
	DB *DBConnection
}

// PersonRepository defines the method available from Person Repository
// domain to be used by external layers.
type PersonRepository interface {
    HelloWorld() error
    Add(string) error
    Parent(string, string) error
    Clear() error
}

// Add creates a new Person in the database with label and attribute name.
func (p *Person) Add(name string) error {
    query := fmt.Sprintf("CREATE (%s:Person) SET %s.name = $name;", name, name)
    _, err := p.DB.Session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
        result, err := transaction.Run(
            query,
            map[string]interface{}{"name": name},
        )
        if err != nil {return nil, err}
        if result.Next() {return result.Record().GetByIndex(0), nil}
        return nil, result.Err()
    })
    if err != nil {return err}
	return nil
}

// Parent creates a new property to the received Person.
func (p *Person) Parent(parent, child string) error {
    query := fmt.Sprintf("MATCH (a:Person {name:'%s'}), (b:Person {name:'%s'}) CREATE (a)-[:PARENT]->(b)", parent, child)
    _, err := p.DB.Session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
        result, err := transaction.Run(
            query,
            map[string]interface{}{},
        )
        if err != nil {return nil, err}
        if result.Next() {return result.Record().GetByIndex(0), nil}
        return nil, result.Err()
    })
    if err != nil {return err}
	return nil
}

// Clear removes all nodes and relationships from the database.
func (p *Person) Clear() error {
    query := `MATCH (n) DETACH DELETE n`
    _, err := p.DB.Session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
        result, err := transaction.Run(
            query,
            map[string]interface{}{},
        )
        if err != nil {return nil, err}
        if result.Next() {return result.Record().GetByIndex(0), nil}
        return nil, result.Err()
    })
    if err != nil {return err}
	return nil
}

func (p *Person) HelloWorld() error {
	greeting, err := p.DB.Session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
        result, err := transaction.Run(
            "CREATE (a:Greeting) SET a.message = $message RETURN a.message + ', from node ' + id(a)",
            map[string]interface{}{"message": "hello, world"})
        if err != nil {
            return nil, err
        }

        if result.Next() {
            return result.Record().GetByIndex(0), nil
        }

        return nil, result.Err()
    })
    if err != nil {
        return err
	}
	
	fmt.Println(greeting.(string))
	return nil
}