package main

import (
	"fmt"
	"net/http"

	"github.com/trongquan1297/cc_score/config"
	"github.com/trongquan1297/cc_score/database"
	"github.com/trongquan1297/cc_score/scoreboard"
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
