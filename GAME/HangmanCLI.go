package Hangman

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
)

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
	MessageR       string
	//MessageI       string
}

var (
	mu               sync.Mutex
	GameDato         GameData
	Mots             []string
	Lettresproposees = make(map[string]bool) //verification
)

func SetWord(level string) {
	switch level {
	case "Facile":
		GameDato.WordToGuess = ChoisirMot(Mots, 2, 4)
	case "Moyen":
		GameDato.WordToGuess = ChoisirMot(Mots, 5, 6)
	case "Difficile":
		GameDato.WordToGuess = ChoisirMot(Mots, 7, 11)
	case "Goldlevel":
		GameDato.WordToGuess = ChoisirMot(Mots, 12, 30)
	default:
		GameDato.WordToGuess = ChoisirMot(Mots, 2, 700)
	}

}

func Contains(slice []string, str string) bool {
	for _, k := range slice {
		if k == str {
			return true
		}
	}
	return false
}

// Mise à jour de la vie
func UpdateLife(lettre string) {
	if !strings.Contains(GameDato.WordToGuess, lettre) {
		GameDato.Life--
	} else {
		GameDato.Score += 2
	}
}

func CheckWin(wordToGuess string, guessedLetters []string) bool {
	for _, letter := range wordToGuess {
		if !Contains(guessedLetters, string(letter)) {
			return false
		}
	}
	return true
}

// choix du mot en fonction de sa longueur
// pour l'attribuer à la  difficulté correspondante
func ChoisirMot(mots []string, minLen, maxLen int) string {
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
	GameDato.Affichage = make([]string, len(GameDato.WordToGuess))
	for i := range GameDato.WordToGuess {
		GameDato.Affichage[i] = "_ "
	}
}
func AfficherMot(guess string) {
	wordRunes := []rune(GameDato.WordToGuess)
	guessRunes := []rune(guess)[0]
	for i, char := range wordRunes {
		if char == guessRunes {
			GameDato.Affichage[i] = string(guess)
		}
	}
}

func LettreDejaProposee(lettre string, Lettresproposees map[string]bool) bool {
	for k := range Lettresproposees {
		if strings.ContainsRune(k, []rune(lettre)[0]) {
			return true
		}
	}
	Lettresproposees[lettre] = true
	return false
}

// pour  les mots du dictionnaires de facon aléatoires
func RetreiveWord() {
	content, err := os.ReadFile("DICTIONNAIRE/mots.txt")
	if err != nil {
		log.Fatal(err)
	}
	Mots = strings.Split(string(content), "\n")
}

func SaveScore(name string, score int) {
	mu.Lock()
	defer mu.Unlock()

	// Ouvrir le fichier en mode append
	file, err := os.OpenFile("GAME/score.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier score.txt :", err)
		return
	}
	defer file.Close()

	// Écrire le score dans le fichier
	_, err = fmt.Fprintf(file, "Nom: %s, Score: %d\n", name, score)
	if err != nil {
		fmt.Println("Erreur lors de l'écriture dans le fichier score.txt :", err)
	}
}
