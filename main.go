package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type todoItem struct {
	Title string
	Done  bool
}

func main() {
	todos := map[int]*todoItem{
		0: {
			Title: "Membaca",
			Done:  false,
		},
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.ParseFiles("app/todo.html"))

		t.ExecuteTemplate(w, "todo.html", map[string]interface{}{"Todos": todos})
	})

	mux.HandleFunc("/proses-todo", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		newTodo := r.PostForm.Get("todo")

		todos[len(todos)] = &todoItem{
			Title: newTodo,
			Done:  false,
		}
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	})

	mux.HandleFunc("/done-todo", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		todoIndex := r.URL.Query().Get("index")

		todoIndexInt, _ := strconv.Atoi(todoIndex)

		if todo, ok := todos[todoIndexInt]; ok {
			todo.Done = !todo.Done

			var responseMessage string
			if todo.Done {
				responseMessage = fmt.Sprintf("%s is now marked as done 😊", todo.Title)
			} else {
				responseMessage = fmt.Sprintf("%s is now marked as undone 😊", todo.Title)
			}
			w.Write([]byte(responseMessage))
		} else {
			http.Error(w, "Todo item not found", http.StatusNotFound)
			return
		}
	})

	server := http.Server{
		Addr:    "localhost:5000",
		Handler: mux,
	}

	errors := server.ListenAndServe()

	if errors != nil {
		panic(errors)
	}
}
