package main

import (
	//Hangman "hangman/GAME"
	routeur "hangman/routeur"
	initTemplate "hangman/templates"
)

func main() {
	initTemplate.InitTemplate()

	routeur.InitServe()
}
