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
	Life            int        //vie restante
	Name            []DataUser //nom du joueur
	IsWin           bool       //indique le joueur a trouvé toutes les letttres
	IsLost          bool       //indique que le joueur n'a pas trouvé toutes les lettres

}
type DataUser struct {
	Name string
}

var (
	counter  int
	temp     *template.Template
	err      error
	user     DataUser
	gameData GameData
	mots     []string
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	temp.ExecuteTemplate(w, "result", nil)
}
func resultTemplate(w http.ResponseWriter, r *http.Request) {
	temp.ExecuteTemplate(w, "home", nil)
}

func playHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		temp.ExecuteTemplate(w, "play", nil)

	} else if r.Method == "POST" {
		guessedLetter := r.FormValue("guessedLetter")
		//mise à jour
		gameData.Life = UpdateLife(gameData.Life, guessedLetter)
		gameData.GuessedLetters = append(gameData.GuessedLetters, guessedLetter)
		gameData.ProposedLetters[guessedLetter] = true
		gameData.IsWin = checkWin(gameData.WordToGuess, gameData.GuessedLetters)
		gameData.IsLost = checkLost(gameData.Life)

		if gameData.IsWin || gameData.IsLost {
			http.Redirect(w, r, "/result", http.StatusSeeOther)
		} else {
			temp.ExecuteTemplate(w, "play", gameData)
		}
	}
	Difficulty := r.FormValue("Difficulty")
	switch Difficulty {
	case "Facile":
		http.Redirect(w, r, "game/easy", http.StatusSeeOther)
	case "Moyen":
		http.Redirect(w, r, "game/medium", http.StatusSeeOther)
	case "Difficile":
		http.Redirect(w, r, "game/hard", http.StatusSeeOther)
	case "Gold Level":
		http.Redirect(w, r, "game/hard level", http.StatusSeeOther)
	default:
		http.Redirect(w, r, "game/easy", http.StatusSeeOther)
	}
}

func setWord(level GameData) {

	switch level.Difficulty {
	case "Facile":
		level.WordToGuess = choisirMot(mots, 2, 6)
	case "Moyen":
		level.WordToGuess = choisirMot(mots, 7, 10)
	case "Difficile":
		level.WordToGuess = choisirMot(mots, 11, 18)
	case "Gold Level":
		level.WordToGuess = choisirMot(mots, 19, 100)
	default:
		level.WordToGuess = choisirMot(mots, 2, 700)
	}
	gameData = level
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
	file, err := os.Open("mots.txt")
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

	setWord(gameData)

	fileServer := http.FileServer(http.Dir("CSS"))
	http.Handle("/CSS/", http.StripPrefix("/CSS/", fileServer))
	http.HandleFunc("/home", HomeHandler)
	http.HandleFunc("/result", HomeHandler)
	http.ListenAndServe(port, nil)

	//

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
