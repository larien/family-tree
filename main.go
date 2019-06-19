package main

import (
	"fmt"
	"log"
	"github.com/larien/family-tree/repository"
)

func main(){
	fmt.Println("Hello, Mundipagg!")

	repository.New()
	log.Println("Repository layer created")

}