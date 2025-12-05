package main

import (
	"cc_score/pkg/config"
	"cc_score/pkg/database"
	"cc_score/pkg/scoreboard"
	"cc_score/pkg/viewer"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

	mux := http.NewServeMux()

	mux.HandleFunc("/viewpost", viewer.ViewPostHandler)

	mux.HandleFunc("/getTopPlayers", func(w http.ResponseWriter, r *http.Request) {
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

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		message := map[string]string{"message": "Scoreboard API!"}
		json.NewEncoder(w).Encode(message)
	})

	mux.HandleFunc("/addScore", scoreboard.PostAddScoreHandler(scoreBoard))
	mux.HandleFunc("/getHighestScore", scoreboard.GetHighestScoreHandler(scoreBoard))
	mux.HandleFunc("/getTop5Players", scoreBoard.GetTopPlayersHandler())
	mux.HandleFunc("/getpostviews", viewer.GetPostViewsHandler)

	fmt.Println("Server is running on :8080")
	if err := http.ListenAndServe(":8080", corsMiddleware(mux)); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if isAllowedOrigin(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
			w.Header().Set("Vary", "Origin")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isAllowedOrigin(origin string) bool {
	if origin == "" {
		return false
	}

	parsed, err := url.Parse(origin)
	if err != nil {
		return false
	}

	host := parsed.Hostname()
	return host == "quannguyen.online" || host == "score.quannguyen.online"
}
