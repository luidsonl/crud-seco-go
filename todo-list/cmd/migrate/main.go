package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"todo-list/internal/database/migration"

	_ "modernc.org/sqlite"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: migrate [up|down|status]")
	}

	db, err := sql.Open("sqlite", "data/app.db")

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
