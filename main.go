package main

import (
	"fmt"
	"log"
	"github.com/larien/family-tree/repository"
	"github.com/larien/family-tree/controller"
)

func main(){
	fmt.Println("Hello, Mundipagg!")

	r := repository.New()
	log.Println("Repository layer created")

	controller.New(r)
	log.Println("Controller layer created")
}