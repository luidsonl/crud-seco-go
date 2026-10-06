package migration

import (
	"database/sql"
	"fmt"
)

func Up(db *sql.DB, migrations []Migration) error {
	if err := CreateTable(db); err != nil {
		return err
	}

	currentVersion, err := CurrentVersion(db)

	if err != nil {
		return err
	}

	for _, migration := range migrations {
		if migration.Version <= currentVersion {
			continue
		}

		tx, err := db.Begin()
		if err != nil {
			return err
		}

		_, err = tx.Exec(migration.Up)
		if err != nil {
			tx.Rollback()
			return err
		}

		if err := SetVersion(tx, migration.Version); err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

func Down(db *sql.DB, migrations []Migration) error {
	if err := CreateTable(db); err != nil {
		return err
	}

	currentVersion, err := CurrentVersion(db)

	if err != nil {
		return err
	}

	for _, migration := range migrations {
		if migration.Version != currentVersion {
			continue
		}

		tx, err := db.Begin()
		if err != nil {
			return err
		}

		_, err = tx.Exec(migration.Down)
		if err != nil {
			tx.Rollback()
			return err
		}

		if err := RemoveVersion(tx, migration.Version); err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}

		break
	}

	return nil
}

func Status(db *sql.DB, migrations []Migration) error {
	if err := CreateTable(db); err != nil {
		return err
	}

	applied, err := AppliedVersions(db)

	if err != nil {
		return err
	}

	appliedMap := make(map[int]bool)

	for _, version := range applied {
		appliedMap[version] = true
	}

	for _, migration := range migrations {
		status := "pending"

		if appliedMap[migration.Version] {
			status = "applied"
		}

		fmt.Printf(
			"%03d_%s [%s]\n",
			migration.Version,
			migration.Name,
			status,
		)
	}

	if len(migrations) < 1 {
		fmt.Println("No migrations found")
	}

	return nil

}
