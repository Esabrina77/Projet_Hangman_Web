package hangman

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
	Affichage      []string
	WORD           string
	Gameletters    string
}

var (
	temp             *template.Template
	err              error
	gameData         GameData
	mots             []string
	lettresproposees = make(map[string]bool) //verification
)

func initHandler(w http.ResponseWriter, r *http.Request) {

	gameData.Life = 4
	gameData.IsWin = checkWin(gameData.WordToGuess, gameData.GuessedLetters)

	gameData.Difficulty = strings.TrimPrefix(r.URL.Path, "/init/")
	switch gameData.Difficulty {
	case "Facile", "Moyen", "Difficile", "Goldlevel":
		break
	default:
		http.Redirect(w, r, "/selection", http.StatusMovedPermanently)
	}
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
	gameData.WORD = strings.Join(gameData.Affichage, "")
	gameData.Gameletters = strings.Join(gameData.GuessedLetters, ", ")

	fmt.Println(gameData.Affichage)
	fmt.Println(gameData.WordToGuess)
	temp.ExecuteTemplate(w, "play", gameData)
}

func getOutHandler(w http.ResponseWriter, r *http.Request) {
	temp.ExecuteTemplate(w, "getOut", nil)
}
func resultHandler(w http.ResponseWriter, r *http.Request) {
	resultData := gameData
	gameData = GameData{} //reset le jeu
	gameData.Name = resultData.Name
	lettresproposees = make(map[string]bool)
	temp.ExecuteTemplate(w, "result", resultData)
}

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
func UpdateLife(lettre string) {
	if !strings.Contains(gameData.WordToGuess, lettre) {
		gameData.Life--
	} else {
		gameData.Score++
	}
}

func checkWin(wordToGuess string, guessedLetters []string) bool {
	for _, letter := range wordToGuess {
		if !contains(guessedLetters, string(letter)) {
			return false
		}
	}
	return true
}

func guessHandler(w http.ResponseWriter, r *http.Request) {
	lettre := r.FormValue("guessedLetter")
	fmt.Println(lettre)

	var toutesLesLettresTrouvees bool

	if !lettreDejaProposee(lettre, lettresproposees) {
		afficherMot(lettre)
		UpdateLife(lettre)
		fmt.Println("wordtoguess: ", gameData.WordToGuess, "  gameData.Affichage: ", strings.Join(gameData.Affichage, ""))
		if gameData.WordToGuess == strings.Join(gameData.Affichage, "") {
			toutesLesLettresTrouvees = true
			fmt.Println("Check win: ", gameData.WordToGuess == strings.Join(gameData.Affichage, ""))
		}
		fmt.Println("Life: ", gameData.Life, "  win: ", toutesLesLettresTrouvees)
		if gameData.Life == 0 || toutesLesLettresTrouvees {
			http.Redirect(w, r, "/result", http.StatusSeeOther)
		}
	}
	gameData.GuessedLetters = append(gameData.GuessedLetters, lettre)
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
			gameData.Affichage[i] = string(guess)
		}
	}
}

func lettreDejaProposee(lettre string, lettresproposees map[string]bool) bool {
	for k := range lettresproposees {
		if strings.ContainsRune(k, []rune(lettre)[0]) {
			return true
		}
	}
	lettresproposees[lettre] = true
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
func Serveur() {
	temp, err = template.ParseGlob("templates/*.html")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

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
