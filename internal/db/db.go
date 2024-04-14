package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// DBManager управляет базой данных.
type DBManager struct {
	db *sql.DB
}

// NewDBManager создает новый экземпляр DBManager.
func NewDBManager(db *sql.DB) *DBManager {
	return &DBManager{db}
}

// SetupTables создает необходимые таблицы в базе данных.
func (mgr *DBManager) SetupTables() error {
	// Создание таблицы пользователей
	_, err := mgr.db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT,
			password_hash TEXT,
			role TEXT
		)
	`)
	if err != nil {
		return err
	}

	// Другие таблицы и соответствующие запросы...

	return nil
}

func main() {
	// Подключение к базе данных PostgreSQL
	db, err := sql.Open("postgres", "user=youruser dbname=yourdb password=yourpassword sslmode=disable")
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer db.Close()

	// Создание экземпляра DBManager
	mgr := NewDBManager(db)

	// Создание таблиц в базе данных
	err = mgr.SetupTables()
	if err != nil {
		fmt.Println("Error setting up tables:", err)
		return
	}

	fmt.Println("Tables created successfully!")
}
