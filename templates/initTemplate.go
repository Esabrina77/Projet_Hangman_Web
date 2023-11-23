package initTemplate

import (
	"fmt"
	"html/template"
	"os"
)

var (
	Temp *template.Template
)

func InitTemplate() {
	temp, err := template.ParseGlob("templates/*.html")
	if err != nil {
		fmt.Printf("Oupss une erreur lié au Templates : %v", err.Error())
		os.Exit(1)
	}
	Temp = temp
}
