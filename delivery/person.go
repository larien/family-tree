package delivery

import (
	"github.com/gin-gonic/gin"
	"net/http"
	c "github.com/larien/family-tree/controller"
)

// Person contains injected interface from Controller layer.
type Person struct {
	Controller c.PersonController
}

func person(version *gin.RouterGroup, controller c.PersonController){
	log.Println("Person delivery started")
	person := &Person{
		Controller: controller,
	}

	endpoints := version.Group("/person")
	{
		endpoints.GET("", person.findAll)
	}
}

// findAll handles GET /person requests and returns all People.
func (p *Person) findAll(c *gin.Context) {
	people, _ := p.Controller.FindAll()

	c.JSON(
		http.StatusOK,
		people,
	)
}

