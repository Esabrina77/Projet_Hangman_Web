package Hangman

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
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
	Mots             []string
	Lettresproposees = make(map[string]bool) //verification
)

func SetWord(level string, gameData *GameData) {
	switch level {
	case "Facile":
		gameData.WordToGuess = ChoisirMot(Mots, 2, 4)
	case "Moyen":
		gameData.WordToGuess = ChoisirMot(Mots, 5, 6)
	case "Difficile":
		gameData.WordToGuess = ChoisirMot(Mots, 7, 11)
	case "Goldlevel":
		gameData.WordToGuess = ChoisirMot(Mots, 12, 30)
	default:
		gameData.WordToGuess = ChoisirMot(Mots, 2, 700)
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
func UpdateLife(lettre string, gameData *GameData) {
	if !strings.Contains(gameData.WordToGuess, lettre) {
		gameData.Life--
	} else {
		gameData.Score += 2
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
func InitMot(gameData *GameData) {
	gameData.Affichage = make([]string, len(gameData.WordToGuess))
	for i := range gameData.WordToGuess {
		gameData.Affichage[i] = "_ "
	}
}

func AfficherMot(guess string, gameData *GameData) {
	wordRunes := []rune(gameData.WordToGuess)
	guessRunes := []rune(guess)[0]
	for i, char := range wordRunes {
		if char == guessRunes {
			gameData.Affichage[i] = string(guess)
		}
	}
}

func IsLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func LettreDejaProposee(lettre string, Lettresproposees map[string]bool) bool {
	// isLetter vérifie si le caractère est une lettre
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

	currentTime := time.Now()

	// Écrire le score et le nom du joueur dans "score.txt"
	_, err = fmt.Fprintf(file, "Nom: %s, Score: %d Date: %s  \n", name, score, currentTime.Format("2006-01-02 15:04:05"))
	if err != nil {
		fmt.Println("Erreur lors de l'écriture dans le fichier score.txt:", err)
	}
}
