package main

import (
	"fmt"
	"log"
	"github.com/larien/family-tree/repository"
	"github.com/larien/family-tree/controller"
	"github.com/larien/family-tree/delivery"
	"github.com/larien/family-tree/middlewares/config"
)

func main(){
	r, err := repository.New()
	if err != nil {
		log.Fatal(err)
	}
	defer func(){
		r.DB.Session.Close()
		r.DB.Driver.Close()
	}()

	log.Println("Repository layer created")
	
	c := controller.New(r)
	log.Println("Controller layer created")

	router := delivery.New(c)
	log.Println("Delivery layer created")

	router.Run(config.Port)
	log.Printf("Running router on port %s", config.Port)

}
