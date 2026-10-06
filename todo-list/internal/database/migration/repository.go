package migration

import "database/sql"

func CreateTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations(
			version INTEGER PRIMARY KEY
		);
	`)

	return err
}

func CurrentVersion(db *sql.DB) (int, error) {
	var version int

	err := db.QueryRow(`
		SELECT COALESCE(MAX(version), 0)
		FROM schema_migrations;
	`).Scan(&version)

	if err != nil {
		return 0, err
	}

	return version, nil
}

func SetVersion(tx *sql.Tx, version int) error {
	_, err := tx.Exec(`
		INSERT INTO schema_migrations (version)
		VALUES(?); 
	`, version)

	return err
}

func RemoveVersion(tx *sql.Tx, version int) error {
	_, err := tx.Exec(`
		DELETE FROM schema_migrations
		WHERE version = ?
	`, version)

	return err
}

func AppliedVersions(db *sql.DB) ([]int, error) {
	rows, error := db.Query(`
		SELECT version
		FROM schema_migrations
		ORDER BY version ASC
	`)

	if error != nil {
		return nil, error
	}

	defer rows.Close()

	var versions []int

	for rows.Next() {

		var version int

		if err := rows.Scan(&version); err != nil {
			return nil, err
		}

		versions = append(versions, version)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return versions, nil
}
