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

	r := repository.New()
	log.Println("Repository layer created")

	c := controller.New(r)
	log.Println("Controller layer created")

	delivery.New(c)
	log.Println("Delivery layer created")
}