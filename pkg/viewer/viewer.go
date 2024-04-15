package viewer

import (
	"cc_score/pkg/config"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Viewer struct {
	db *sqlx.DB
}

type Post struct {
	ID    string `db:"id"`
	Views int    `db:"views"`
}

func NewViewer(db *sqlx.DB) *Viewer {
	return &Viewer{db: db}
}

func ConnectDB(cfg config.DatabaseConfig) (*sqlx.DB, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	db, err := sqlx.Connect("mysql", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}
	return db, nil
}

func (v *Viewer) ViewPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	id := data.ID
	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	post := Post{}
	err := v.db.Get(&post, "SELECT * FROM posts WHERE id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err := v.db.Exec("INSERT INTO posts (id, views) VALUES (?, 1)", id)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error inserting new post: %v", err), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, fmt.Sprintf("Error querying database: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
		_, err := v.db.Exec("UPDATE posts SET views = views + 1 WHERE id = ?", id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error updating post views: %v", err), http.StatusInternalServerError)
			return
		}
	}

	fmt.Fprintf(w, "Post with ID %d viewed successfully", id)
}

func (v *Viewer) GetViewsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	var views int
	err := v.db.Get(&views, "SELECT views FROM posts WHERE id = ?", id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error querying database: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]int{"views": views}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (v *Viewer) GetPostViewsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	var views int
	err := v.db.Get(&views, "SELECT views FROM posts WHERE id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Error querying database: %v", err), http.StatusInternalServerError)
		return
	}

	response := struct {
		ID    string `json:"id"`
		Views int    `json:"views"`
	}{
		ID:    id,
		Views: views,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshalling JSON: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
