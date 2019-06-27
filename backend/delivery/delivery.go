package delivery

import (
	c "github.com/larien/family-tree/controller"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

// Deliveries contains all Deliveries interface from different domains.
// When creating a new domain, its interface should be placed here in
// order to external layers have access to its methods.
type Deliveries struct {
	Person c.PersonController
}

// New creates a Delivery layer and shares its reference to the external
// layer have access to its methods. It should receive as parameter the
// connections needed by the Delivery layer and the Controllers interface.
// Each Delivery dependency is placed according to each domain.
func New(controllers *c.Controllers) *gin.Engine {
	person := &Deliveries{
		Person: controllers.Person,
	}
	return router(person)
}

// router sets up routing for the github.com/larien/family-tree.
func router(deliveries *Deliveries) *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())
	endpoints(router, deliveries)
	return router
}

// endpoints defines endpoints from each entity from Delivery layer.
func endpoints(router *gin.Engine, deliveries *Deliveries) {
	v1 := router.Group("/api/v1")
	{
		person(v1, deliveries.Person)
	}
}