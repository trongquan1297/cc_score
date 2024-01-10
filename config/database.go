package database

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/your-username/your-project-name/config"
)

// ConnectDB connects to the MySQL database using configuration.
func ConnectDB(cfg config.DatabaseConfig) (*sqlx.DB, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	db, err := sqlx.Connect("mysql", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}
	return db, nil
}
