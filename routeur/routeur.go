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
	http.HandleFunc("/viewscore", controller.ViewScoreHandler)
	
	port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // Fallback au port 8080 si la variable d'environnement PORT n'est pas définie
    }

    if err := http.ListenAndServe("0.0.0.0:"+port, nil); err != nil {
        log.Fatal(err)
    }
}
