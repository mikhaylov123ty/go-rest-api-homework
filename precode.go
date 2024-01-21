package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Структура задачи
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

// Инициализация списка задач с примерами
var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postman",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Обработчик запросов GET для эндпоинта /users
func getTasks(w http.ResponseWriter, r *http.Request) {

	resp, err := json.MarshalIndent(tasks, "", "    ")
	if err != nil {
		fmt.Printf("Ошибка в процессе сериализации задач: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// Обработчик запросов POST для эндпоинта /users
// Принимает запросы в формате {"id":"3","description":"Отправить финальное задание REST API на ревью","note":"Закончил. Ура!","applications":["GoLand","git"]}
func postTasks(w http.ResponseWriter, r *http.Request) {
	task := &Task{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Ошибка в процессе чтения тела запроса задач: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(body, &task); err != nil {
		fmt.Printf("Ошибка в процессе десериализации задачи: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := tasks[task.ID]; ok {
		fmt.Printf("Задача с номером %s уже существует\n", task.ID)
		http.Error(w, "Задача с таким номером уже существует", http.StatusBadRequest)
		return
	}

	tasks[task.ID] = *task

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Задача успешно добавлена\n"))
}

// Обработчик запросов GET для эндпоинта /users/id
func getTasksID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if _, ok := tasks[id]; !ok {
		fmt.Printf("Задачи с ID #%s не существует\n", id)
		http.Error(w, "Задачи с таким ID не существует", http.StatusBadRequest)
		return
	}

	resp, err := json.MarshalIndent(tasks[id], "", "    ")
	if err != nil {
		fmt.Printf("Ошибка в процессе сериализации задачи: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// Обработчик запросов DELETE для эндпоинта /users/id
func deleteTasksID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if _, ok := tasks[id]; !ok {
		fmt.Printf("Задачи с ID #%s не существует\n", id)
		http.Error(w, "Задачи с таким ID не существует", http.StatusBadRequest)
		return
	}

	delete(tasks, id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Задача успешно удалена\n"))
}

func main() {
	r := chi.NewRouter()
	r.Get("/tasks", getTasks)
	r.Post("/tasks", postTasks)
	r.Get("/tasks/{id}", getTasksID)
	r.Delete("/tasks/{id}", deleteTasksID)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
