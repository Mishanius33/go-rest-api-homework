package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

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
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта
// Обработчик для получения всех задач
func getTasks(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// Обработчик для отправки задачи на сервер
func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &newTask); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tasks[newTask.ID] = newTask
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// Обработчик для получения задачи по ID
func getTaskById(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	task, id := tasks[taskID]
	if !id {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Обработчик удаления задачи по ID
func deleteTaskById(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")

	w.Header().Set("Content-Type", "application/json")
	if _, id := tasks[taskID]; !id {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	delete(tasks, taskID)
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	r.Get("/tasks", getTasks)
	r.Post("/tasks", createTask)
	r.Get("/tasks/{id}", getTaskById)
	r.Delete("/tasks/{id}", deleteTaskById)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
