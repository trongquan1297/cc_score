package main

import (
	"cc_score/pkg/config"
	"cc_score/pkg/database"
	"cc_score/pkg/scoreboard"
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

	scoreBoard := scoreboard.NewScoreBoard(db)

	http.HandleFunc("/addScore", scoreBoard.AddScoreHandler)
	http.HandleFunc("/getHighestScore", scoreBoard.GetHighestScoreHandler)

	fmt.Println("Server is running on :8080")
	http.ListenAndServe(":8080", nil)
}
