package db

import (
	"fmt"
	"strconv"
)

type Task struct {
	ID      string `json:"id" db:"id"`
	Date    string `json:"date" db:"date"`
	Title   string `json:"title" db:"title"`
	Comment string `json:"comment" db:"comment"`
	Repeat  string `json:"repeat" db:"repeat"`
}

func Tasks(limit int) ([]*Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler 
	          ORDER BY date ASC, id ASC 
	          LIMIT ?`

	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения задач: %v", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var task Task
		var id int64
		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования задачи: %v", err)
		}
		task.ID = strconv.FormatInt(id, 10)
		tasks = append(tasks, &task)
	}

	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}
