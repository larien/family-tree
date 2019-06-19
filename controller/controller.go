package controller

import r "github.com/larien/family-tree/repository"

// Controllers contains all Controllers interface from different domains.
// When creating a new domain, its interface should be placed here in
// order to external layers have access to its methods.
type Controllers struct {}


// New creates a Controllers layer and shares its reference to the external
// layer have access to its methods. It should receive as parameter the
// connections needed by the Controllers layer and the Repository interface.
// Each Controllers dependency is placed according to each domain.
func New(repository *r.Repositories) *Controllers {
	return &Controllers{
		// implement controllers domains here
	}
}
