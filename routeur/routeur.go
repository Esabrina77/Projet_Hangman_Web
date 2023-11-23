package routeur

import (
	"hangman/controller"
	"log"
	"net/http"
)

func InitServe() {
	fileServer := http.FileServer(http.Dir("CSS"))
	http.Handle("/CSS/", http.StripPrefix("/CSS/", fileServer))
	http.HandleFunc("/home", controller.HomeHandler)
	http.HandleFunc("/result", controller.ResultHandler)
	http.HandleFunc("/selection", controller.SelectionHandler)
	http.HandleFunc("/init/", controller.InitHandler)
	http.HandleFunc("/play", controller.PlayHandler)
	http.HandleFunc("/getOut", controller.GetOutHandler)
	http.HandleFunc("/help", controller.HelpHandler)
	http.HandleFunc("/treatment", controller.TreatHandler)
	http.HandleFunc("/guess", controller.GuessHandler)
	if err := http.ListenAndServe(controller.Port, nil); err != nil {
		log.Fatal(err)
	}
}
