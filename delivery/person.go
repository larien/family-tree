package delivery

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
	"github.com/larien/family-tree/entity"
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
		endpoints.POST("", person.add)
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

// add handles POST /person requests and adds People and its relationships.
func (p *Person) add(c *gin.Context) {
	log.Println("Add")
	var people []entity.Person

	if err := c.BindJSON(&people); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": "Failed to parse json",
				"error":   err,
			})
		return
	}

	err := p.Controller.Add(people)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed register people",
				"error":   err,
			})
		return
	}

	c.JSON(
		http.StatusCreated,
		gin.H{
			"status":  http.StatusCreated,
			"message": "People registered successfully!",
		})
}
