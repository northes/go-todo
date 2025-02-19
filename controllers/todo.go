package controllers

import (
 "github.com/gorilla/mux"
 "github.com/ichtrojan/go-todo/config"
 "github.com/ichtrojan/go-todo/models"
 "github.com/rs/zerolog/log"
 "html/template"
 "net/http"
)

var (
	id        int
	item      string
	completed int
	view      = template.Must(template.ParseFiles("./views/index.html"))
	database  = config.Database()
)

func Show(w http.ResponseWriter, r *http.Request) {
 statement, err := database.Query(`SELECT * FROM todos`)

 if err != nil {
  log.Error().Err(err).Str("operation", "show_todos").Msg("Failed to query todos")
  http.Error(w, "Internal Server Error", http.StatusInternalServerError)
  return
 }

 var todos []models.Todo

	for statement.Next() {
  err = statement.Scan(&id, &item, &completed)

  if err != nil {
   log.Error().Err(err).Str("operation", "scan_todo").Int("todo_id", id).Msg("Failed to scan todo")
   continue
  }

  todo := models.Todo{
			Id:        id,
			Item:      item,
			Completed: completed,
		}

		todos = append(todos, todo)
	}

	data := models.View{
		Todos: todos,
	}

	_ = view.Execute(w, data)
}

func Add(w http.ResponseWriter, r *http.Request) {
 item := r.FormValue("item")

 _, err := database.Exec(`INSERT INTO todos (item) VALUE (?)`, item)

 if err != nil {
  log.Error().Err(err).Str("operation", "add_todo").Str("item", item).Msg("Failed to add todo")
  http.Error(w, "Internal Server Error", http.StatusInternalServerError)
  return
 }

 log.Info().Str("operation", "add_todo").Str("item", item).Msg("Todo added successfully")

	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
 vars := mux.Vars(r)
 id := vars["id"]

 _, err := database.Exec(`DELETE FROM todos WHERE id = ?`, id)

 if err != nil {
  log.Error().Err(err).Str("operation", "delete_todo").Str("todo_id", id).Msg("Failed to delete todo")
  http.Error(w, "Internal Server Error", http.StatusInternalServerError)
  return
 }

 log.Info().Str("operation", "delete_todo").Str("todo_id", id).Msg("Todo deleted successfully")

	http.Redirect(w, r, "/", 301)
}

func Complete(w http.ResponseWriter, r *http.Request) {
 vars := mux.Vars(r)
 id := vars["id"]

 _, err := database.Exec(`UPDATE todos SET completed = 1 WHERE id = ?`, id)

 if err != nil {
  log.Error().Err(err).Str("operation", "complete_todo").Str("todo_id", id).Msg("Failed to complete todo")
  http.Error(w, "Internal Server Error", http.StatusInternalServerError)
  return
 }

 log.Info().Str("operation", "complete_todo").Str("todo_id", id).Msg("Todo marked as complete")

	http.Redirect(w, r, "/", 301)
}
