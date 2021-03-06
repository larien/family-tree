package delivery

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
	"fmt"
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
		endpoints.GET("/name/:name", person.find)
		endpoints.GET("/ascendancy/:name", person.ascendancy)
		endpoints.POST("", person.add)
		endpoints.GET("/clear", person.clear)
	}
}

// findAll handles GET /person requests and returns all People.
func (p *Person) findAll(c *gin.Context) {
	people, err := p.Controller.FindAll()
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "Failed to find all people",
				"error":   err,
			})
		return
	}

	if people == nil {
		c.JSON(
			http.StatusNoContent,
			gin.H{
				"message": "No people were found",
			})
		return
	}

	c.JSON(
		http.StatusOK,
		people,
	)
}


// find handles GET /person/name/:name requests and return the Person.
func (p *Person) find(c *gin.Context) {
	name := c.Param("name")
	
	person, err := p.Controller.Find(name)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": fmt.Sprintf("Failed to find %s", name),
				"error":   err,
			})
		return
	}

	if person == nil {
		c.JSON(
			http.StatusNoContent,
			gin.H{
				"message": fmt.Sprintf("%s wasn't found", name),
			})
		return
	}

	c.JSON(
		http.StatusOK,
		person,
	)
}

// add handles POST /person requests and adds People and its relationships.
func (p *Person) add(c *gin.Context) {
	var people []entity.Person

	if err := c.BindJSON(&people); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
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
				"message": "Failed to register people",
				"error":   err,
			})
		return
	}

	c.JSON(
		http.StatusCreated,
		gin.H{
			"message": "People registered successfully!",
		})
}

// find handles GET /ascendancy/:name requests and return the Person's
// ascendancy.
func (p *Person) ascendancy(c *gin.Context) {
	name := c.Param("name")

	ascendants, err := p.Controller.Ascendancy(name)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": fmt.Sprintf("Failed to find %s's ascendancy", name),
				"error":   err,
			})
		return
	}

	if ascendants == nil {
		c.JSON(
			http.StatusNoContent,
			gin.H{
				"message": fmt.Sprintf("%s wasn't found", name),
			})
		return
	}

	c.JSON(
		http.StatusOK,
		ascendants,
	)
}

// clear is a helper function that cleans the database.
func (p *Person) clear(c *gin.Context) {
	err := p.Controller.Clear()
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": fmt.Sprintf("Failed to clear database"),
				"error":   err,
			})
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "Database cleared sucessfully",
		})
}