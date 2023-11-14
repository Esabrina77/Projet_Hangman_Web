package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

// VARIABLES & CONSTANTES
const port = "localhost:8080"

// structure pour le stock des données de la partie en cours
type GameData struct {
	WordToGuess     string   //mot à deviner
	Difficulty      string   // pour choisir le niveau de difficulté
	GuessedLetters  []string //lettre devinées
	ProposedLetters map[string]bool
	Score           int
	Life            int    //vie restante
	Name            string //nom du joueur
	IsWin           bool   //indique le joueur a trouvé toutes les letttres
	IsLost          bool   //indique que le joueur n'a pas trouvé toutes les lettres

}

var (
	counter  int
	temp     *template.Template
	err      error
	gameData GameData
	mots     []string
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	temp.ExecuteTemplate(w, "home", nil)
}

func TreatHandler(w http.ResponseWriter, r *http.Request) {
	gameData.Name = r.FormValue("name")

	if gameData.Name == "" {
		errorMessage := "VEILLEZ REMPLIR TOUS  LES CHAMPS DU FORMULAIRE"
		http.Redirect(w, r, "/user/home?error="+errorMessage, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/selection", http.StatusSeeOther)

}

func selectionHandler(w http.ResponseWriter, r *http.Request) {
	temp.ExecuteTemplate(w, "selection", gameData.Name)
}

func easyHandler(w http.ResponseWriter, r *http.Request) {
	setWord("Facile")
	data := GameData{
		Name:        gameData.Name,
		Life:        gameData.Life,
		WordToGuess: gameData.WordToGuess,
	}

	temp.ExecuteTemplate(w, "easy", data)
}
func mediumHandler(w http.ResponseWriter, r *http.Request) {
	setWord("moyen")
	data := GameData{
		Name:        gameData.Name,
		Life:        gameData.Life,
		WordToGuess: gameData.WordToGuess,
	}
	temp.ExecuteTemplate(w, "medium", data)
}
func hardHandler(w http.ResponseWriter, r *http.Request) {
	setWord("Difficile")
	data := GameData{
		Name:        gameData.Name,
		Life:        gameData.Life,
		WordToGuess: gameData.WordToGuess,
	}
	temp.ExecuteTemplate(w, "hard", data)
}
func goldlevelHandler(w http.ResponseWriter, r *http.Request) {
	setWord("Goldlevel")
	data := GameData{
		Name:        gameData.Name,
		Life:        gameData.Life,
		WordToGuess: gameData.WordToGuess,
	}
	temp.ExecuteTemplate(w, "goldlevel", data)
}

func getOutHandler(w http.ResponseWriter, r *http.Request) {
	temp.ExecuteTemplate(w, "getOut", nil)
}
func resultHandler(w http.ResponseWriter, r *http.Request) {
	temp.ExecuteTemplate(w, "result", gameData)
}

func setWord(level string) {
	switch level {
	case "Facile":
		gameData.WordToGuess = choisirMot(mots, 2, 6)
	case "Moyen":
		gameData.WordToGuess = choisirMot(mots, 6, 11)
	case "Difficile":
		gameData.WordToGuess = choisirMot(mots, 11, 19)
	case "Goldlevel":
		gameData.WordToGuess = choisirMot(mots, 19, 100)
	default:
		gameData.WordToGuess = choisirMot(mots, 2, 700)
	}

}

func contains(slice []string, str string) bool {
	for _, k := range slice {
		if k == str {
			return true
		}
	}
	return false
}

func checkWin(wordToGuess string, guessedLetters []string) bool {
	for _, letter := range wordToGuess {
		if !contains(guessedLetters, string(letter)) {
			return false
		}
	}
	return true
}

func checkLost(life int) bool {
	return life <= 0
}

func main() {
	//Ouverture du fichier
	file, err := os.Open("DICTIONNAIRE/mots.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Lecture des mots du dictionnaire
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		mots = append(mots, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	temp, err = template.ParseGlob("./templates/*.html")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	bookServer := http.FileServer(http.Dir("DICTIONNAIRE"))
	http.Handle("/DICTIONNAIRE/", http.StripPrefix("/DICTIONNAIRE/", bookServer))

	fileServer := http.FileServer(http.Dir("CSS"))
	http.Handle("/CSS/", http.StripPrefix("/CSS/", fileServer))
	http.HandleFunc("/home", homeHandler)
	http.HandleFunc("/result", resultHandler)
	http.HandleFunc("/selection", selectionHandler)
	http.HandleFunc("/easy", easyHandler)
	http.HandleFunc("/medium", mediumHandler)
	http.HandleFunc("/hard", hardHandler)
	http.HandleFunc("/goldlevel", goldlevelHandler)
	http.HandleFunc("/getOut", getOutHandler)
	http.HandleFunc("/treatment", TreatHandler)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}

}

// choix du mot pour la difficulté
func choisirMot(mots []string, minLen, maxLen int) string {
	var motsFiltres []string
	for _, mot := range mots {
		l := len(mot)
		if l >= minLen && l <= maxLen {
			motsFiltres = append(motsFiltres, mot)
		}
	}
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(motsFiltres))
	return motsFiltres[index]
}

// Mise à jour de la vie
func UpdateLife(life int, letter string) int {
	if !strings.Contains(gameData.WordToGuess, letter) {
		life--
	}
	return life
}

func afficherPendu(pendu []string, vie int) {
	if vie < len(pendu) {
		fmt.Println(pendu[vie])
	} else {
		fmt.Println(pendu[len(pendu)-1])
	}
}

func splitWordToGuess(mot string) []string {
	display := make([]string, len(mot))
	for i, l := range mot {
		display[i] = string(l)
	}
	return display
}
