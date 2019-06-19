package repository

// Repositories contains all Repositories interface from different domains.
// When creating a new domain, its interface should be placed here in
// order to external layers have access to its methods.
type Repositories struct {}


// New creates a Repository layer and shares its reference to the external
// layer have access to its methods. It should receive as parameter
// the connections needed by the Repository layer. Each Repository dependency
// is placed according to each domain.
func New() *Repositories {
	return &Repositories{
		// implement repository domains here
	}
}
