package main

import (
	"cc_score/pkg/config"
	"cc_score/pkg/database"
	"cc_score/pkg/scoreboard"
	"cc_score/pkg/viewer"
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {

	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	db, err := database.ConnectDB(cfg.Database)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer db.Close()

	viewer := viewer.NewViewer(db)

	scoreBoard := scoreboard.NewScoreBoard(db)

	http.HandleFunc("/viewpost", viewer.ViewPostHandler)

	http.HandleFunc("/getTopPlayers", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		topPlayers, err := scoreBoard.GetTopPlayers()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(topPlayers)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		message := map[string]string{"message": "Scoreboard API!"}
		json.NewEncoder(w).Encode(message)
	})

	http.HandleFunc("/addScore", scoreboard.PostAddScoreHandler(scoreBoard))
	http.HandleFunc("/getHighestScore", scoreboard.GetHighestScoreHandler(scoreBoard))
	http.HandleFunc("/getTop5Players", scoreBoard.GetTopPlayersHandler())
	http.HandleFunc("/getpostviews", viewer.GetPostViewsHandler)
	fmt.Println("Server is running on :8080")
	http.ListenAndServe(":8080", nil)
}
