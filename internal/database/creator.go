package database

import (
	"os"

	"github.com/jmoiron/sqlx"
)

func ReadSQLFile(filePath string) ([]byte, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func ExecuteSQLCommands(db *sqlx.DB, commands []byte) error {
	_, err := db.Exec(string(commands))
	return err
}
