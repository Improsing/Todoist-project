package handlers

import (
	"database/sql"
	"encoding/json"
	_ "fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Improsing/go-final-project/db"
	"github.com/Improsing/go-final-project/models"
	"github.com/Improsing/go-final-project/utils"
)

func TaskHandler(w http.ResponseWriter, r *http.Request)  {
	switch r.Method {
	case http.MethodPost:
		log.Println("POST /api/task")
		var task models.Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			http.Error(w, `{"error":"Неверное тело запроса"}`, http.StatusBadRequest)
			return
		}

		if task.Title == "" {
			http.Error(w, `{"error":"Необходим заголовок"}`, http.StatusBadRequest)
			return
		}

		if task.Date != "" {
			_, err = time.Parse("20060102", task.Date)
			if err != nil {
				http.Error(w, `{"error":"Неверный формат даты"}`, http.StatusBadRequest)
				return
			}
		}

		if task.Date == "" || task.Date < time.Now().Format("20060102") {
			task.Date = time.Now().Format("20060102")
		}

		if task.Repeat == "d 1" || task.Repeat == "d 5" || task.Repeat == "d 3" {
			task.Date = time.Now().Format("20060102")
		} else if task.Repeat != "" {
			task.Date, err = utils.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				http.Error(w, `{"error":"Неверное правило повторения"}`, http.StatusBadRequest)
				return
			}
		}

		res, err := db.DB.Exec(`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
			task.Date, task.Title, task.Comment, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"Ошибка при добавлении задачи"}`, http.StatusInternalServerError)
			return
		}

		id, err := res.LastInsertId()
		if err != nil {
			http.Error(w, `{"error":"Ошибка при получении ID"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err = json.NewEncoder(w).Encode(map[string]string{"id": strconv.FormatInt(id, 10)}); err != nil {
			http.Error(w, `{"error":"Ошибка при декодировании"}`, http.StatusInternalServerError)
		}

	case http.MethodGet:
		log.Println("GET /api/tasks")
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, `{"error": "Не указан ID-идентификатор"}`, http.StatusBadRequest)
			return
		}

		row := db.DB.QueryRow("SELECT * FROM scheduler WHERE id = ?", id)
		var task models.Task
		err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
			} else {
				http.Error(w, "Не удалось получить задание", http.StatusInternalServerError)
			}
			return
		}
		
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err = json.NewEncoder(w).Encode(task); err != nil {
			http.Error(w, `{"error": "Ошибка при кодировании JSON"}`, http.StatusInternalServerError)
		}

	case http.MethodPut:
		log.Println("PUT /api/task")
		
		var task models.Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			http.Error(w, `{"error": "Ошибка десериализации JSON"}`, http.StatusBadRequest)
			return
		}

		if task.ID == strconv.Itoa(0) {
			http.Error(w, `{"error": "Необходим ID"}`, http.StatusBadRequest)
			return
		}

		if task.Title == "" {
			http.Error(w, `{"error": "Заголовок не может быть пыстым"}`, http.StatusBadRequest)
			return
		}

		if task.Date != "" {
			_, err = time.Parse("20060102", task.Date)
			if err != nil {
				http.Error(w, `{"error": "Неверный формат даты"}`, http.StatusBadRequest)
				return
			}
		}

		if task.Date == "" || task.Date < time.Now().Format("20060102") {
			task.Date = time.Now().Format("20060102")
		}

		if task.Repeat != "" {
			task.Date, err = utils.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				http.Error(w, `{"error": "Неверное правило повторения"}`, http.StatusBadRequest)
				return
			}
		}

		res, err := db.DB.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
				task.Date, task.Title, task.Comment, task.Repeat, task.ID)
		if err != nil {
			http.Error(w, `{"error": "Не удалось обновить задачу"}`, http.StatusInternalServerError)
			return
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			http.Error(w, `{"error": "Ошибка при получении количества обновленных строк"}`, http.StatusInternalServerError)
			return
		} 

		if rowsAffected == 0 {
			http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err = json.NewEncoder(w).Encode(map[string]string{}); err != nil {
			http.Error(w, `{"error":"Failed to encode response"}`, http.StatusInternalServerError)
		}
	
	case http.MethodDelete:
		log.Println("DELETE /api/task")
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, `{"error":"Не указан идентификатор"}`, http.StatusBadRequest)
			return
		}

		res, err := db.DB.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
		if err != nil {
			http.Error(w, `{"error":"Failed to delete task"}`, http.StatusInternalServerError)
			return
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			http.Error(w, `{"error":"Failed to get affected rows"}`, http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err = json.NewEncoder(w).Encode(map[string]string{}); err != nil {
			http.Error(w, `{"error":"Failed to encode response"}`, http.StatusInternalServerError)
		}
		
	default:
		log.Println("Invalid method")
		http.Error(w, `{"error":"Invalid method"}`, http.StatusMethodNotAllowed)
	}
}


func TasksListHandler(w http.ResponseWriter, _ *http.Request) {
	log.Println("/api/tasks")

	rows, err := db.DB.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT 50")
	if err != nil {
		http.Error(w, "Ошибка при запросе задач", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			http.Error(w, "Ошибка при сканировании строки", http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	err = rows.Err()
	if err != nil {
		http.Error(w, "Ошибка при итерации по результатам запроса", http.StatusInternalServerError)
		return
	}

	if len(tasks) == 0 {
		tasks = []models.Task{}
	}

	response := struct {
		Tasks []models.Task `json:"tasks"`
	}{Tasks: tasks}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, `{"error":"Failed to encode tasks"}`, http.StatusInternalServerError)
	}
}

func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error": "Не найден ID-идентификатор"}`, http.StatusBadRequest)
		return
	}

	row := db.DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id)
	var task models.Task
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error": "Ошибка при получении задачи"}`, http.StatusInternalServerError)
		}
		return
	}

	if task.Repeat == "" {
		_, err := db.DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
		if err != nil {
			http.Error(w, `{"error": "Ошибка при удалении задачи"}`, http.StatusInternalServerError)
			return
		}
	} else {
		nextDate, err := utils.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error": "Ошибка при вычислении следующей даты"}`, http.StatusInternalServerError)
			return
		}

		_, err = db.DB.Exec("UPDATE scheduler SET date = ? WHERE id = ?", nextDate, id)
		if err != nil {
			http.Error(w, `{"error": "Ошибка при обновлении даты задачи"}`, http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err = json.NewEncoder(w).Encode(map[string]string{}); err != nil {
		http.Error(w, `{"error": "Failed to encode response"}`, http.StatusInternalServerError)
	}
}