package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"sync"
)

var (
	db   *sqlx.DB
	once sync.Once
)

// InitializeDatabase initializes the database connection.
func InitializeDatabase(dbName, dbPort, dbHost, dbUser, dbPassword, sslmode string) {
	var connStr string
	if dbPassword == "" {
		connStr = fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=%s",
			dbHost, dbUser, dbName, dbPort, sslmode)
	} else {
		connStr = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			dbHost, dbUser, dbPassword, dbName, dbPort, sslmode)
	}

	var err error
	once.Do(func() {
		db, err = sqlx.Connect("postgres", connStr)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
	})

	log.Println("Database connection established successfully.")
}

// GetDB returns the database connection.
func GetDB() *sqlx.DB {
	return db
}

// CloseDB closes the database connection.
func CloseDB() {
	if db != nil {
		err := db.Close()
		if err != nil {
			log.Printf("Error closing database: %v", err)
		} else {
			log.Println("Database connection closed.")
		}
	}
}

func SetupDatabase() {
	// Create the uuid-ossp extension
	if _, err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";"); err != nil {
		log.Fatalf("Failed to create uuid-ossp extension: %v", err)
	}

	createTables()
}

func createTables() {
	userInfoTable := `
	CREATE TABLE IF NOT EXISTS user_info(
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), 
		first_name VARCHAR(50) NOT NULL,
		last_name VARCHAR(50) NOT NULL,
	    college_name VARCHAR(100),
	    college_state VARCHAR(50),
		account_create_date TIMESTAMPTZ NOT NULL,  -- TIMESTAMP WITH TIME ZONE
		leetcode_username VARCHAR(50),
		codechef_username VARCHAR(50),
		codeforces_username VARCHAR(50),
		min_duration_in_seconds INT NOT NULL,
		max_duration_in_seconds INT NOT NULL
	);`

	siteTable := `
	CREATE TABLE IF NOT EXISTS user_site_info(
		id SERIAL PRIMARY KEY,
		user_id UUID REFERENCES user_info(id) ON DELETE CASCADE,
		site_name TEXT NOT NULL ,
		is_site_enabled BOOLEAN NOT NULL,
		is_automatic_calendar_notification_enabled BOOLEAN NOT NULL,
		seconds_before_which_app_notification_to_set INTEGER[],
		CONSTRAINT unique_user_site UNIQUE (user_id, site_name)
	);`

	// Execute the queries
	_, err := db.Exec(userInfoTable)
	if err != nil {
		log.Fatalf("Error creating user_info table: %v", err)
	}

	_, err = db.Exec(siteTable)
	if err != nil {
		log.Fatalf("Error creating site table: %v", err)
	}

	log.Println("All tables created successfully.")
}
