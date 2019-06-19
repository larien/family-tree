package delivery

import c "github.com/larien/family-tree/controller"

// Deliveries contains all Deliveries interface from different domains.
// When creating a new domain, its interface should be placed here in
// order to external layers have access to its methods.
type Deliveries struct {}


// New creates a Delivery layer and shares its reference to the external
// layer have access to its methods. It should receive as parameter the
// connections needed by the Delivery layer and the Controllers interface.
// Each Delivery dependency is placed according to each domain.
func New(controller *c.Controllers){}
