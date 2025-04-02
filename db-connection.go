package main

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/marcboeker/go-duckdb" // Import DuckDB driver
)

var (
	db   *sql.DB
	once sync.Once
)

// GetDBInstance returns a singleton DuckDB connection with a persistent database file
func GetDBInstance() *sql.DB {
	once.Do(func() {
		var err error
		// Store DuckDB database in a directory
		db, err = sql.Open("duckdb", "database/store.duckdb") // Use a file instead of in-memory
		if err != nil {
			log.Fatal("Failed to connect to DuckDB:", err)
		}
		_, dbErr := db.Exec("CREATE TABLE IF NOT EXISTS webmap (id UUID, url TEXT, path TEXT, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (id))")
		if dbErr != nil {
			log.Fatal("Failed to create table:", err)
		}
		log.Println("DuckDB connected successfully (Persistent Mode)")
	})
	return db
}

// func main() {
// 	// Get the singleton DuckDB instance
// 	database := GetDBInstance()
// 	defer database.Close()

// 	// Example: Create a table
// 	_, err := database.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER, name TEXT)")
// 	if err != nil {
// 		log.Fatal("Failed to create table:", err)
// 	}

// 	// Insert data
// 	_, err = database.Exec("INSERT INTO users (id, name) VALUES (1, 'Alice')")
// 	if err != nil {
// 		log.Fatal("Failed to insert data:", err)
// 	}

// 	fmt.Println("Data inserted successfully")

// 	// Query data
// 	rows, err := database.Query("SELECT id, name FROM users")
// 	if err != nil {
// 		log.Fatal("Failed to query data:", err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var id int
// 		var name string
// 		if err := rows.Scan(&id, &name); err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Printf("User: %d - %s\n", id, name)
// 	}
// }
