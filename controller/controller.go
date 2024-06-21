package controller

import (
	"fmt"
	Hangman "hangman/GAME"
	initTemplate "hangman/templates"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

// VARIABLES & CONSTANTES
var (
	mu        sync.Mutex
	Activated bool
)

// initialisation de l'espace de jeu
func InitHandler(w http.ResponseWriter, r *http.Request) {
	Activated = true
	Hangman.GameDato.Life = 9
	Hangman.GameDato.IsWin = Hangman.CheckWin(Hangman.GameDato.WordToGuess, Hangman.GameDato.GuessedLetters)

	//strings.TrimPrefix permet de
	Hangman.GameDato.Difficulty = strings.TrimPrefix(r.URL.Path, "/init/")
	switch Hangman.GameDato.Difficulty {
	case "Facile", "Moyen", "Difficile", "Goldlevel":
		break
	default:
		http.Redirect(w, r, "/selection", http.StatusMovedPermanently)
	}
	Hangman.RetreiveWord()
	Hangman.SetWord(Hangman.GameDato.Difficulty)
	Hangman.InitMot()

	http.Redirect(w, r, "/play", http.StatusSeeOther)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	initTemplate.Temp.ExecuteTemplate(w, "home", nil)
}

func TreatHandler(w http.ResponseWriter, r *http.Request) {
	Hangman.GameDato.Name = r.FormValue("name")
	if Hangman.GameDato.Name == "" {
		errorMessage := "VEILLEZ REMPLIR TOUS  LES CHAMPS DU FORMULAIRE"
		http.Redirect(w, r, "/user/home?error="+errorMessage, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/selection", http.StatusSeeOther)

}

func SelectionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	initTemplate.Temp.ExecuteTemplate(w, "selection", Hangman.GameDato.Name)
}

func PlayHandler(w http.ResponseWriter, r *http.Request) {
	// if !Activated {
	// 	http.Redirect(w, r, "/selection", http.StatusSeeOther)
	// 	return
	// }
	Hangman.GameDato.WORD = strings.Join(Hangman.GameDato.Affichage, "")
	Hangman.GameDato.Gameletters = strings.Join(Hangman.GameDato.GuessedLetters, ", ")

	fmt.Println(Hangman.GameDato.Affichage)
	fmt.Println(Hangman.GameDato.WordToGuess)
	initTemplate.Temp.ExecuteTemplate(w, "play", Hangman.GameDato)
}

func GetOutHandler(w http.ResponseWriter, r *http.Request) {
	initTemplate.Temp.ExecuteTemplate(w, "getOut", nil)
}

func HelpHandler(w http.ResponseWriter, r *http.Request) {
	initTemplate.Temp.ExecuteTemplate(w, "Help", nil)
}

func ResultHandler(w http.ResponseWriter, r *http.Request) {
	resultData := Hangman.GameDato

	Hangman.GameDato = Hangman.GameData{} //reset le jeu
	Hangman.GameDato.Name = resultData.Name
	Hangman.Lettresproposees = make(map[string]bool)
	Hangman.SaveScore(resultData.Name, resultData.Score)
	//SAUVEGARDE DU SCORE DES JOUEURS
	Activated = false

	initTemplate.Temp.ExecuteTemplate(w, "result", resultData)
}
func GuessHandler(w http.ResponseWriter, r *http.Request) {
	lettre := r.FormValue("guessedLetter")
	lettreByte := lettre[0]
	_ = Hangman.IsLetter(lettreByte)

	fmt.Println(lettreByte)
	var toutesLesLettresTrouvees bool

	if !Hangman.LettreDejaProposee(lettre, Hangman.Lettresproposees) {
		Hangman.GameDato.MessageR = ""
		Hangman.AfficherMot(lettre)
		Hangman.UpdateLife(lettre)
		fmt.Println("wordtoguess: ", Hangman.GameDato.WordToGuess, "  Hangman.GameDato.Affichage: ", strings.Join(Hangman.GameDato.Affichage, ""))
		if Hangman.GameDato.WordToGuess == strings.Join(Hangman.GameDato.Affichage, "") {
			toutesLesLettresTrouvees = true
			fmt.Println("Check win: ", Hangman.GameDato.WordToGuess == strings.Join(Hangman.GameDato.Affichage, ""))
		}
		fmt.Println("Life: ", Hangman.GameDato.Life, "  win: ", toutesLesLettresTrouvees)
		if Hangman.GameDato.Life == 0 || toutesLesLettresTrouvees {
			http.Redirect(w, r, "/result", http.StatusSeeOther)
			return
		}
		//Hangman.GameDato.MessageI = "LETTRE INCORRECT"
	} else {
		//Hangman.GameDato.MessageI = ""
		Hangman.GameDato.MessageR = "lettre déjà proposée"
		fmt.Println("lettre déjà proposée")
	}

	Hangman.GameDato.GuessedLetters = append(Hangman.GameDato.GuessedLetters, lettre)
	http.Redirect(w, r, "/play", http.StatusSeeOther)

}

func ViewScoreHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	// Lire le contenu du fichier score.txt
	content, err := ioutil.ReadFile("GAME/score.txt")
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de la lecture du fichier score.txt : %v", err), http.StatusInternalServerError)
		return
	}
	// content, err := os.ReadFile("GAME/score.txt")
	// if err != nil {
	// 	http.Error(w, fmt.Sprintf("Erreur lors de la lecture du fichier score.txt : %v", err), http.StatusInternalServerError)
	// 	return
	// }

	// Envoyer le contenu comme réponse HTTP
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}
