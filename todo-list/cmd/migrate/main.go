package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"todo-list/internal/database/migration"

	_ "modernc.org/sqlite"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: migrate [up|down|status]")
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dbDir := filepath.Join(wd, "data")
	dbPath := filepath.Join(dbDir, "app.db")

	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite", dbPath)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	migrations, err := migration.Load("./migrations")

	if err != nil {
		log.Fatal(err)
	}

	switch os.Args[1] {
	case "up":
		err = migration.Up(db, migrations)

	case "down":
		err = migration.Down(db, migrations)

	case "status":
		err = migration.Status(db, migrations)

	default:
		log.Fatalf("Unknown command: %s", os.Args[1])
	}

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done")

}
