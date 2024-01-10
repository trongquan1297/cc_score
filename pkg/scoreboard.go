package scoreboard

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/jmoiron/sqlx"
)

// Player struct represents the player information.
type Player struct {
	ID    string `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Score int    `json:"score" db:"score"`
}

// ScoreBoard struct represents the scoreboard.
type ScoreBoard struct {
	mu sync.Mutex
	db *sqlx.DB
}

// NewScoreBoard creates a new ScoreBoard instance.
func NewScoreBoard(db *sqlx.DB) *ScoreBoard {
	return &ScoreBoard{
		db: db,
	}
}

// AddPlayer adds a new player to the scoreboard.
func (sb *ScoreBoard) AddPlayer(player Player) error {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	_, err := sb.db.NamedExec(`INSERT INTO players (name, score) VALUES (:name, :score)`, player)
	return err
}

// GetHighestScore returns the player with the highest score.
func (sb *ScoreBoard) GetHighestScore() (Player, error) {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	var highestScorePlayer Player
	err := sb.db.Get(&highestScorePlayer, `SELECT * FROM players ORDER BY score DESC LIMIT 1`)
	if err == sql.ErrNoRows {
		return Player{}, nil
	}

	return highestScorePlayer, err
}

// AddScoreHandler handles the API endpoint for adding score.
func (sb *ScoreBoard) AddScoreHandler(w http.ResponseWriter, r *http.Request) {
	var newPlayer Player
	err := json.NewDecoder(r.Body).Decode(&newPlayer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = sb.AddPlayer(newPlayer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetHighestScoreHandler handles the API endpoint for getting the highest score.
func (sb *ScoreBoard) GetHighestScoreHandler(w http.ResponseWriter, r *http.Request) {
	highestScorePlayer, err := sb.GetHighestScore()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if highestScorePlayer.ID == "" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(highestScorePlayer)
}
