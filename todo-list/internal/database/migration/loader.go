package migration

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
)

var migrationPattern = regexp.MustCompile(
	`^(\d+)_(.+)\.(up|down)\.sql$`,
)

func Load(dir string) ([]Migration, error) {
	entries, err := os.ReadDir(dir)

	if err != nil {
		return nil, err
	}

	migrations := make(map[int]*Migration)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		match := migrationPattern.FindStringSubmatch(entry.Name())

		if match == nil {
			continue
		}

		content, err := os.ReadFile(filepath.Join(dir, entry.Name()))

		if err != nil {
			return nil, err
		}

		version, err := strconv.Atoi(match[1])

		if err != nil {
			return nil, err
		}

		name := match[2]

		direction := match[3]

		if _, exists := migrations[version]; !exists {
			migrations[version] = &Migration{
				Version: version,
				Name:    name,
			}
		}

		if direction == "up" {
			migrations[version].Up = string(content)
		} else {
			migrations[version].Down = string(content)
		}
	}

	result := []Migration{}

	for _, migration := range migrations {
		result = append(result, *migration)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Version < result[j].Version
	})

	return result, nil

}
