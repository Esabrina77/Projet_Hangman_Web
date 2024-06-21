package controller

import (
	"fmt"
	Hangman "hangman/GAME"
	initTemplate "hangman/templates"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/sessions"
)

// VARIABLES & CONSTANTES
var (
	mu        sync.Mutex
	Activated bool
	store     = sessions.NewCookieStore([]byte("secret-key"))
)

func getSessionData(r *http.Request) (*Hangman.GameData, *sessions.Session, error) {
	session, err := store.Get(r, "hangman-session")
	if err != nil {
		return nil, nil, err
	}
	gameData, ok := session.Values["gameData"].(*Hangman.GameData)
	if !ok || gameData == nil {
		gameData = &Hangman.GameData{}
		session.Values["gameData"] = gameData
	}
	return gameData, session, nil
}

// initialisation de l'espace de jeu
func InitHandler(w http.ResponseWriter, r *http.Request) {
	Activated = true
	gameData, session, err := getSessionData(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	gameData.Life = 9
	gameData.IsWin = Hangman.CheckWin(gameData.WordToGuess, gameData.GuessedLetters)

	//strings.TrimPrefix permet de
	gameData.Difficulty = strings.TrimPrefix(r.URL.Path, "/init/")
	switch gameData.Difficulty {
	case "Facile", "Moyen", "Difficile", "Goldlevel":
		break
	default:
		http.Redirect(w, r, "/selection", http.StatusMovedPermanently)
		return
	}
	Hangman.RetreiveWord()
	Hangman.SetWord(gameData.Difficulty, gameData)
	Hangman.InitMot(gameData)
	if err = session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/play", http.StatusSeeOther)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	initTemplate.Temp.ExecuteTemplate(w, "home", nil)
}

func TreatHandler(w http.ResponseWriter, r *http.Request) {
	gameData, session, err := getSessionData(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	gameData.Name = r.FormValue("name")
	if gameData.Name == "" {
		errorMessage := "VEILLEZ REMPLIR TOUS LES CHAMPS DU FORMULAIRE"
		http.Redirect(w, r, "/user/home?error="+errorMessage, http.StatusSeeOther)
		return
	}
	if err = session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/selection", http.StatusSeeOther)
}

func SelectionHandler(w http.ResponseWriter, r *http.Request) {
	gameData, _, err := getSessionData(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	initTemplate.Temp.ExecuteTemplate(w, "selection", gameData.Name)
}

func PlayHandler(w http.ResponseWriter, r *http.Request) {
	gameData, _, err := getSessionData(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// if !Activated {
	// 	http.Redirect(w, r, "/selection", http.StatusSeeOther)
	// 	return
	// }
	gameData.WORD = strings.Join(gameData.Affichage, "")
	gameData.Gameletters = strings.Join(gameData.GuessedLetters, ", ")

	fmt.Println(gameData.Affichage)
	fmt.Println(gameData.WordToGuess)
	initTemplate.Temp.ExecuteTemplate(w, "play", gameData)
}

func GetOutHandler(w http.ResponseWriter, r *http.Request) {
	initTemplate.Temp.ExecuteTemplate(w, "getOut", nil)
}

func HelpHandler(w http.ResponseWriter, r *http.Request) {
	initTemplate.Temp.ExecuteTemplate(w, "Help", nil)
}

func ResultHandler(w http.ResponseWriter, r *http.Request) {
	gameData, session, err := getSessionData(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resultData := *gameData

	*gameData = Hangman.GameData{} //reset le jeu
	gameData.Name = resultData.Name
	Hangman.Lettresproposees = make(map[string]bool)
	Hangman.SaveScore(resultData.Name, resultData.Score)
	//SAUVEGARDE DU SCORE DES JOUEURS
	Activated = false

	if err = session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	initTemplate.Temp.ExecuteTemplate(w, "result", resultData)
}

func GuessHandler(w http.ResponseWriter, r *http.Request) {
	gameData, session, err := getSessionData(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lettre := r.FormValue("guessedLetter")
	lettreByte := lettre[0]
	_ = Hangman.IsLetter(lettreByte)

	fmt.Println(lettreByte)
	var toutesLesLettresTrouvees bool

	if !Hangman.LettreDejaProposee(lettre, Hangman.Lettresproposees) {
		gameData.MessageR = ""
		Hangman.AfficherMot(lettre, gameData)
		Hangman.UpdateLife(lettre, gameData)
		fmt.Println("wordtoguess: ", gameData.WordToGuess, "  gameData.Affichage: ", strings.Join(gameData.Affichage, ""))
		if gameData.WordToGuess == strings.Join(gameData.Affichage, "") {
			toutesLesLettresTrouvees = true
			fmt.Println("Check win: ", gameData.WordToGuess == strings.Join(gameData.Affichage, ""))
		}
		fmt.Println("Life: ", gameData.Life, "  win: ", toutesLesLettresTrouvees)
		if gameData.Life == 0 || toutesLesLettresTrouvees {
			if err = session.Save(r, w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/result", http.StatusSeeOther)
			return
		}
		//gameData.MessageI = "LETTRE INCORRECT"
	} else {
		//gameData.MessageI = ""
		gameData.MessageR = "lettre déjà proposée"
		fmt.Println("lettre déjà proposée")
	}

	gameData.GuessedLetters = append(gameData.GuessedLetters, lettre)
	if err = session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

	// Envoyer le contenu comme réponse HTTP
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write(content)

}
