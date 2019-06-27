package repository

import (
    "github.com/larien/family-tree/entity"
    "github.com/neo4j/neo4j-go-driver/neo4j"
    "encoding/json"
    "io/ioutil"
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
    Add(string) error
    Backup(string) error
    Children(string) ([]string, error)
    Clear() error
    Connected(string) ([]string, error)
    DeleteWithoutChildren() error
    Parent(string, string) error
    Retrieve(string) (*entity.Person, error)
    RetrieveAll() ([]entity.Person, error)
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

// Backup realizes a backup from the current data in the database.
func (p *Person) Backup(filename string) error {
	people, err := p.RetrieveAll()
	if err != nil {return err}
	
	return dump(people, filename)
}

// Children returns the current Person's children.
func (p *Person) Children(name string) ([]string, error) {
    query := fmt.Sprintf("MATCH (n:Person)-[:PARENT]->(m) WHERE n.name = '%s' RETURN DISTINCT collect(m.name) as names;", name)
    peopleNames, err := p.DB.Session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
        result, err := transaction.Run(
            query,
            map[string]interface{}{})
        if err != nil {
            return nil, err
        }
        var peopleNames []string
        for result.Next() {
            record := result.Record()

            names, ok := record.Get("names")
            if !ok {return nil, fmt.Errorf("Couldn't get names")}

            var n []string
            if names != nil {
                n, err = parseInterfaceToString(names)
                if err != nil {return nil, err}
            }

            peopleNames = n
        }
        return peopleNames, result.Err()
    })
    if err != nil {return nil, err}

    asserted, ok := peopleNames.([]string)
    if !ok {
        return nil, nil
    }

	return asserted, nil
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

// Connected returns all People connected to a Person from the database.
func (p *Person) Connected(name string) ([]string, error) {
    query := fmt.Sprintf("MATCH (u:Person)-[:PARENT*0..]-(connected:Person) WHERE u.name = '%s' RETURN DISTINCT collect(connected.name) as names;", name)
    peopleNames, err := p.DB.Session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
        result, err := transaction.Run(
            query,
            map[string]interface{}{})
        if err != nil {
            return nil, err
        }
        var peopleNames []string
        for result.Next() {
            record := result.Record()

            names, ok := record.Get("names")
            if !ok {return nil, fmt.Errorf("Couldn't get names")}

            var n []string
            if names != nil {
                n, err = parseInterfaceToString(names)
                if err != nil {return nil, err}
            }

            peopleNames = n
        }
        return peopleNames, result.Err()
    })
    if err != nil {return nil, err}

    asserted, ok := peopleNames.([]string)
    if !ok {
        return nil, nil
    }

	return asserted, nil
}

// DeleteWithoutChildren deletes all People without children that have parents.
func (p *Person) DeleteWithoutChildren() error {
    query := fmt.Sprintf("MATCH (a:Person) WHERE not ((a)-[:PARENT]->(:Person)) AND ()-[:PARENT]->(a) DETACH DELETE a;")
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

// Retrieve returns a Person from the database.
func (p *Person) Retrieve(name string) (*entity.Person, error) {
    query := fmt.Sprintf("MATCH (n: Person {name: '%s'}) RETURN n.name as name, n.parents as parents, n.children as children", name)
    person, err := p.DB.Session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
        result, err := transaction.Run(
            query,
            map[string]interface{}{})
        if err != nil {
            return nil, err
        }
        for result.Next() {
            record := result.Record()
            
            name, ok := record.Get("name")
            if !ok {return nil, fmt.Errorf("Couldn't get name")}
            parents, ok := record.Get("parents")
            if !ok {return nil, fmt.Errorf("Couldn't get children")}
            children, ok := record.Get("children")
            if !ok {return nil, fmt.Errorf("Couldn't get children")}
 
            var p []string
            if parents != nil {
                p, err = parseInterfaceToString(parents)
                if err != nil {return nil, err}
            }
            
            var c []string
            if children != nil {
                c, err = parseInterfaceToString(children)
                if err != nil {return nil, err}
            }

            person := &entity.Person{
                Name: name.(string),
                Parents: p,
                Children: c,
            }
            return person, nil
        }
        return nil, result.Err()
    })
    if err != nil {return nil, err}

    asserted, ok := person.(*entity.Person)
    if !ok {
        return nil, nil
    }
	return asserted, nil
}

// RetrieveAll returns all People from the database.
func (p *Person) RetrieveAll() ([]entity.Person, error) {
    query := fmt.Sprintf("MATCH (n) RETURN n.name as name, n.parents as parents, n.children as children")
    people, err := p.DB.Session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
        result, err := transaction.Run(
            query,
            map[string]interface{}{})
        if err != nil {
            return nil, err
        }
        people := []entity.Person{}
        for result.Next() {
            record := result.Record()

            name, ok := record.Get("name")
            if !ok {return nil, fmt.Errorf("Couldn't get name")}
            parents, ok := record.Get("parents")
            if !ok {return nil, fmt.Errorf("Couldn't get children")}
            children, ok := record.Get("children")
            if !ok {return nil, fmt.Errorf("Couldn't get children")}

            var p []string
            if parents != nil {
            p, err = parseInterfaceToString(parents)
            if err != nil {return nil, err}
            }

            var c []string
            if children != nil {
                c, err = parseInterfaceToString(children)
                if err != nil {return nil, err}
            }

            person := entity.Person{
                Name: name.(string),
                Parents: p,
                Children: c,
            }
            people = append(people, person)
        }
        return people, result.Err()
    })
    if err != nil {return nil, err}

    asserted, ok := people.([]entity.Person)
    if !ok {
        return nil, nil
    }

	return asserted, nil
}

// dump creates a backup file to save current data.
func dump(people []entity.Person, filename string) error {
	peopleJSON, err := json.Marshal(people)
	if err != nil {return err}

	err = ioutil.WriteFile(filename, peopleJSON, 0644)
	if err != nil {return err}

	log.Printf("Dump saved to %s", filename)

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