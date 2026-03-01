package db

import (
	"database/sql"
	"final-project/pkg/date"
	"fmt"
	"strconv"
	"time"
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

func GetTask(id string) (*Task, error) {
	var task Task
	var dbID int64

	taskID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("некорректный идентификатор задачи")
	}

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err = DB.QueryRow(query, taskID).Scan(&dbID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("задача не найдена")
		}
		return nil, fmt.Errorf("ошибка получения задачи: %v", err)
	}

	task.ID = strconv.FormatInt(dbID, 10)
	return &task, nil
}

func UpdateTask(task *Task) error {
	taskID, err := strconv.ParseInt(task.ID, 10, 64)
	if err != nil {
		return fmt.Errorf("некорректный идентификатор задачи")
	}

	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	result, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, taskID)
	if err != nil {
		return fmt.Errorf("ошибка обновления задачи: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка проверки обновления: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

func DeleteTask(id string) error {
	taskID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("некорректный идентификатор задачи")
	}

	query := `DELETE FROM scheduler WHERE id = ?`
	result, err := DB.Exec(query, taskID)
	if err != nil {
		return fmt.Errorf("ошибка удаления задачи: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка проверки удаления: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

func CompleteTask(id string) error {
	task, err := GetTask(id)
	if err != nil {
		return err
	}

	if task.Repeat != "" {
		now := time.Now()
		nextDate, err := date.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
		task.Date = nextDate
		return UpdateTask(task)
	} else {
		return DeleteTask(id)
	}
}

func UpdateDate(next string, id string) error {
	taskID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("некорректный идентификатор задачи")
	}

	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	result, err := DB.Exec(query, next, taskID)
	if err != nil {
		return fmt.Errorf("ошибка обновления даты: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка проверки обновления: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}
