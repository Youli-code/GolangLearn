package main

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func initDB() *sql.DB {
	db, err := sql.Open("sqlite", "./tasks.db")
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}

	schema := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		completed BOOLEAN,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(schema); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	ensureUpdatedAtColumn(db)

	return db
}

func ensureUpdatedAtColumn(db *sql.DB) {
	rows, err := db.Query(`PRAGMA table_info(tasks)`)
	if err != nil {
		log.Printf("schema check failed: %v", err)
		return
	}
	defer rows.Close()

	hasUpdatedAt := false
	for rows.Next() {
		var cid int
		var name, colType string
		var dfltVal sql.NullString
		var notNull, pk int
		if err := rows.Scan(&cid, &name, &colType, &notNull, &dfltVal, &pk); err != nil {
			log.Printf("schema scan failed: %v", err)
			return
		}
		if name == "updated_at" {
			hasUpdatedAt = true
			break
		}
	}
	if err := rows.Err(); err != nil {
		log.Printf("schema rows err: %v", err)
		return
	}

	if !hasUpdatedAt {
		if _, err := db.Exec(`ALTER TABLE tasks ADD COLUMN updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP`); err != nil {
			log.Printf("add updated_at failed (maybe already exists?): %v", err)
		} else {
			log.Printf("migrated: added tasks.updated_at")
		}
	}
}
