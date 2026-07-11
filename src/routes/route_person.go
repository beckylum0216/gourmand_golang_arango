package routes

import (
	"github.com/gin-gonic/gin"
	"gourmand.golang.arango/src/entities"
	"gourmand.golang.arango/src/interfaces"
)

type PersonRoute struct {
	service interfaces.IPerson
}

func NewPersonRoutes(service interfaces.IPerson) *PersonRoute {
	return &PersonRoute{
		service: service,
	}
}

func (pr *PersonRoute) CreatePerson(c *gin.Context) {
	ctx := c.Request.Context()
	var person entities.Person

	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	if _, err := pr.service.CreatePerson(ctx, &person); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"message": "Person created successfully",
		"person":  person,
	})
}

func (pr *PersonRoute) GetPerson(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	person, err := pr.service.GetPerson(ctx, id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Person not found"})
	}

	c.JSON(200, person)
}

func (pr *PersonRoute) GetPersons(c *gin.Context) {
	ctx := c.Request.Context()
	persons, err := pr.service.GetPersons(ctx)
	if err != nil {
		c.JSON(404, gin.H{"error": "No persons found"})
		return
	}
	c.JSON(200, persons)
}

func (pr *PersonRoute) UpdatePerson(c *gin.Context) {
	ctx := c.Request.Context()
	var personData = &entities.Person{}
	if err := c.ShouldBindJSON(personData); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}
	id := c.Param("id")
	err := pr.service.UpdatePerson(ctx, id, personData)
	if err != nil {
		c.JSON(404, gin.H{"error": "Person not found"})
	}
	c.JSON(200, gin.H{"message": "Person updated successfully"})
}

func (pr *PersonRoute) DeletePerson(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	err := pr.service.DeletePerson(ctx, id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Person not found"})
		return
	}

	c.JSON(200, gin.H{"message": "Person deleted successfully"})
}

func (pr *PersonRoute) GetPersonDetails(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	details, err := pr.service.GetPersonDetails(ctx, id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Person with id " + id + " not found"})
		return
	}
	c.JSON(200, details)
}

func (pr *PersonRoute) GetPeopleWithDetails(c *gin.Context) {
	ctx := c.Request.Context()

	people, err := pr.service.GetPeopleWithDetails(ctx)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, people)
}
