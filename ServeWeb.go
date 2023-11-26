package main

import (
	//Hangman "hangman/GAME"
	"fmt"
	routeur "hangman/routeur"
	initTemplate "hangman/templates"
)

func main() {
	fmt.Println("server is running...")

	fmt.Println("")
	fmt.Print("CLICK HERE to play---> http://localhost:8080/home \n")
	fmt.Println("")
	fmt.Println("TO STOP THE SERVER , PRESS  'ctrl+C' ")
	initTemplate.InitTemplate()
	routeur.InitServe()

}
