package repository

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"log"
)

var(
	dbURI = "bolt://localhost:7687"
	// dbURI = "bolt://localhost:7687"  // TODO - get from environment variables
	dbUser = "neo4j"
	dbPassword = "1234"
)

// Repositories contains all Repositories interface from different domains.
// When creating a new domain, its interface should be placed here in
// order to external layers have access to its methods.
type Repositories struct {
	DB *DBConnection
	Person PersonRepository
}

type DBConnection struct {
	Driver neo4j.Driver
	Session neo4j.Session
}

// New creates a Repository layer and shares its reference to the external
// layer have access to its methods. It should receive as parameter
// the connections needed by the Repository layer. Each Repository dependency
// is placed according to each domain.
func New() (*Repositories, error) {
	connection, err := dbConnection(dbURI, dbUser, dbPassword)
	if err != nil {
        return nil, err
    }
	log.Println("Database connection stablished")

	return &Repositories{
		DB: connection,
		Person: newPersonRepository(connection),
		// implement repository domains here
	}, nil
}

func dbConnection(uri, username, password string) (*DBConnection, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
    if err != nil {
        return nil, err
    }

    session, err := driver.Session(neo4j.AccessModeWrite)
    if err != nil {
        return nil, err
	}
	
	return &DBConnection{
		Driver: driver,
		Session: session,
	}, nil
}