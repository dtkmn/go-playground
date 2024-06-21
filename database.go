package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	// Connect to the database
	connStr := os.Getenv("DATABASE_URL")

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to the database")
}

// Resource represents a resource in the database
type Resource struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GetResources fetches all resources from the database
func GetResources() ([]Resource, error) {
	rows, err := db.Query("SELECT id, name FROM resources")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resources := []Resource{}
	for rows.Next() {
		var resource Resource
		if err := rows.Scan(&resource.ID, &resource.Name); err != nil {
			return nil, err
		}
		resources = append(resources, resource)
	}
	return resources, nil
}

// AddResource adds a new resource to the database
func AddResource(name string) error {
	_, err := db.Exec("INSERT INTO resources (name) VALUES ($1)", name)
	return err
}
