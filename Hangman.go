package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

// VARIABLES & CONSTANTES
const port = "localhost:8080"

// structure pour le stock des données de la partie en cours
type GameData struct {
	WordToGuess    string   //mot à deviner
	Difficulty     string   // pour choisir le niveau de difficulté
	GuessedLetters []string //lettre devinées
	Score          int
	Life           int    //vie restante
	Name           string //nom du joueur
	IsWin          bool   //indique le joueur a trouvé toutes les letttres
	IsLost         bool   //indique que le joueur n'a pas trouvé toutes les lettres
	Affichage      []string
}

var (
	pendu            []string
	lettresDevinees  = make([]bool, len(gameData.WordToGuess))
	temp             *template.Template
	err              error
	gameData         GameData
	mots             []string
	lettresproposees = make(map[string]bool) //verification
)

func initHandler(w http.ResponseWriter, r *http.Request) {
	gameData.Difficulty = strings.TrimPrefix(r.URL.Path, "/init/")
	gameData.Life = 9
	retreiveWord()
	setWord(gameData.Difficulty)
	InitMot()

	http.Redirect(w, r, "/play", http.StatusSeeOther)
}

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

func playHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println(gameData.Affichage)
	fmt.Println(gameData.WordToGuess)
	temp.ExecuteTemplate(w, "play", gameData)
}

func getOutHandler(w http.ResponseWriter, r *http.Request) {
	temp.ExecuteTemplate(w, "getOut", nil)
}
func resultHandler(w http.ResponseWriter, r *http.Request) {
	temp.ExecuteTemplate(w, "result", gameData)
}

/*
une liste de mot est attribué à
chaque niveau de difficulté en fonction de la longueur du mot
*/
func setWord(level string) {
	switch level {
	case "Facile":
		gameData.WordToGuess = choisirMot(mots, 2, 4)
	case "Moyen":
		gameData.WordToGuess = choisirMot(mots, 5, 6)
	case "Difficile":
		gameData.WordToGuess = choisirMot(mots, 7, 11)
	case "Goldlevel":
		gameData.WordToGuess = choisirMot(mots, 12, 30)
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

// Mise à jour de la vie
func UpdateLife(life int, letter string) int {
	if !strings.Contains(gameData.WordToGuess, letter) {
		life--
	}
	return life
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

func guessHandler(w http.ResponseWriter, r *http.Request) {
	lettre := r.FormValue("guessedLetter")
	fmt.Println(lettre)

	if !lettreDejaProposee(lettre, lettresproposees) {
		afficherMot(lettre)
	}
	http.Redirect(w, r, "/play", http.StatusSeeOther)

}

// choix du mot pour la difficulté
func choisirMot(mots []string, minLen, maxLen int) string {
	var motsFiltres []string
	for _, mot := range mots {
		l := len(mot)
		mot = strings.TrimSpace(mot)
		if l >= minLen && l <= maxLen {
			motsFiltres = append(motsFiltres, mot)
		}
	}

	index := rand.Intn(len(motsFiltres) - 1)
	return motsFiltres[index]
}

// POUR masquer lettres pas encore devinées
func InitMot() {
	gameData.Affichage = make([]string, len(gameData.WordToGuess))
	for i := range gameData.WordToGuess {
		gameData.Affichage[i] = "_ "
	}
}
func afficherMot(guess string) {
	wordRunes := []rune(gameData.WordToGuess)
	guessRunes := []rune(guess)[0]
	for i, char := range wordRunes {
		if char == guessRunes {
			gameData.Affichage[i] = guess
		}
	}
}

func lettreDejaProposee(lettre string, lettresproposees map[string]bool) bool {
	for k := range lettresproposees {
		if strings.ContainsRune(k, []rune(lettre)[0]) {
			return true
		}
	}
	return false
}

// pour  les mots du dictionnaires de facon aléatoires
func retreiveWord() {
	content, err := os.ReadFile("DICTIONNAIRE/mots.txt")
	if err != nil {
		log.Fatal(err)
	}
	mots = strings.Split(string(content), "\n")

}

// FONCTION MAIN---------PRINCIPALE
func main() {
	temp, err = template.ParseGlob("templates/*.html")
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
	http.HandleFunc("/init/", initHandler)
	http.HandleFunc("/play", playHandler)
	http.HandleFunc("/getOut", getOutHandler)
	http.HandleFunc("/treatment", TreatHandler)
	http.HandleFunc("/guess", guessHandler)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}

}
