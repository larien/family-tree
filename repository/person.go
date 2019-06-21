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