package db

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(255) NOT NULL DEFAULT '',
    comment TEXT,
    repeat VARCHAR(128)
);

CREATE INDEX idx_date ON scheduler(date);
`

func Init(dbFile string) error {

	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("ошибка открытия БД: %v", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("ошибка подключения к БД: %v", err)
	}

	if install {
		if _, err := DB.Exec(schema); err != nil {
			return fmt.Errorf("ошибка создания схемы: %v", err)
		}
		fmt.Printf("База данных создана: %s\n", dbFile)
	} else {
		fmt.Printf("База данных подключена: %s\n", dbFile)
	}

	return nil
}

func AddTask(task *Task) (string, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`

	result, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return "", fmt.Errorf("ошибка добавления задачи: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("ошибка получения ID: %v", err)
	}

	return strconv.FormatInt(id, 10), nil
}
