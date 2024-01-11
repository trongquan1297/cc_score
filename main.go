package main

import (
	"cc_score/pkg/config"
	"cc_score/pkg/database"
	"cc_score/pkg/scoreboard"
	"encoding/json"
	"fmt"
	"net/http"
)

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Thiết lập các headers cho CORS
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Cho phép trình duyệt gửi các headers sau khi nhận một phản hồi từ server
		w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

		// Nếu đây là một yêu cầu OPTIONS (pre-flight), không xử lý tiếp và trả về ngay
		if r.Method == "OPTIONS" {
			return
		}

		// Gọi handler tiếp theo trong chuỗi middleware
		next(w, r)
	}
}

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

	scoreBoard := scoreboard.NewScoreBoard(db)

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
		message := map[string]string{"welcome": "Welcome to the Scoreboard API!"}
		json.NewEncoder(w).Encode(message)
	})

	http.HandleFunc("/addScore", scoreboard.PostAddScoreHandler(scoreBoard))
	http.HandleFunc("/getHighestScore", scoreboard.GetHighestScoreHandler(scoreBoard))
	http.HandleFunc("/getTop5Players", corsMiddleware(scoreBoard.GetTopPlayersHandler()))
	fmt.Println("Server is running on :8080")
	http.ListenAndServe(":8080", nil)
}
