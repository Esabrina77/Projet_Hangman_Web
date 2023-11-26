package main

import (
	//Hangman "hangman/GAME"
	"fmt"
	routeur "hangman/routeur"
	initTemplate "hangman/templates"
)

func main() {
	fmt.Println("server is running...")
	fmt.Print("click her to play---> http://localhost:8080/home")
	initTemplate.InitTemplate()
	routeur.InitServe()

}
