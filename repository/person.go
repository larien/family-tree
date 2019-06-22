package repository

import (
    "github.com/larien/family-tree/entity"
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
    Retrieve(string) (*entity.Person, error)
    Add(string) error
    Parent(string, string) error
    HelloWorld() error
    Clear() error
}

// Retrieve returns a Person from the database..
func (p *Person) Retrieve(name string) (*entity.Person, error) {
    query := fmt.Sprintf("MATCH (n: Person {name: '%s'}) RETURN n.parents as parents, n.children as children", name)
    person, err := p.DB.Session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
        result, err := transaction.Run(
            query,
            map[string]interface{}{})
        if err != nil {
            return nil, err
        }
        for result.Next() {
            record := result.Record()
            parents, ok := record.Get("parents")
            if !ok {return nil, fmt.Errorf("Couldn't get children")}
            children, ok := record.Get("children")
            if !ok {return nil, fmt.Errorf("Couldn't get children")}
            
            p, err := parseInterfaceToString(parents)
            if err != nil {return nil, err}

            var c []string
            if children != nil {
                c, err = parseInterfaceToString(children)
                if err != nil {return nil, err}
            }

            person := &entity.Person{
                Name: name,
                Parents: p,
                Children: c,
            }
            return person, nil
        }
        return nil, result.Err()
    })
    if err != nil {return nil, err}
	
	return person.(*entity.Person), nil
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
    // create parent-child relationship
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

    // create parents attribute in child
    query = fmt.Sprintf("MERGE (n: Person {name: '%s'}) SET n.parents = COALESCE(n.parents, []) + '%s'", child, parent)
    _, err = p.DB.Session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
        result, err := transaction.Run(
            query,
            map[string]interface{}{},
        )
        if err != nil {return nil, err}
        if result.Next() {return result.Record().GetByIndex(0), nil}
        return nil, result.Err()
    })
    if err != nil {return err}

    // create children attribute in parent
    query = fmt.Sprintf("MERGE (n: Person {name: '%s'}) SET n.children = COALESCE(n.children, []) + '%s'", parent, child)
    _, err = p.DB.Session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
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

func parseInterfaceToString(i interface{}) (s []string, err error){
    converted, ok := i.([]interface{})
    if !ok {
        return nil, fmt.Errorf("Argument is not a slice")
    }
    for _, v := range converted {
        s = append(s, v.(string))
    }
    return
}