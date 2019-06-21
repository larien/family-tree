package main

import (
	"fmt"
	"log"
	"github.com/larien/family-tree/repository"
	"github.com/larien/family-tree/controller"
	"github.com/larien/family-tree/delivery"
)

func main(){
	fmt.Println("Hello, Mundipagg!")

	r, err := repository.New()
	if err != nil {
		log.Fatal(err)
	}
	r.Person.Add("Lauren")
	
	log.Println("Repository layer created")

	c := controller.New(r)
	log.Println("Controller layer created")

	delivery.New(c)
	log.Println("Delivery layer created")

	r.DB.Session.Close()
	r.DB.Driver.Close()
}