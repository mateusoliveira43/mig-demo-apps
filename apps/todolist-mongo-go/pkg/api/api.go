// MIT License

// Copyright (c) 2020 Mohamad Fadhil

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mig-demo-apps/apps/todolist-mongo-go/pkg/database"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

var toDoTable database.ToDoTable

func CreateItemRoute(w http.ResponseWriter, r *http.Request) {
	description := r.FormValue("description")
	id := toDoTable.CreateItem(description)
	result := make(map[string]string)
	result["Id"] = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func UpdateItemRoute(w http.ResponseWriter, r *http.Request) {
	// Get URL parameter from mux
	vars := mux.Vars(r)
	id := vars["id"]
	completed, _ := strconv.ParseBool(r.FormValue("completed"))
	err, message := toDoTable.UpdateItem(id, completed)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		io.WriteString(w, fmt.Sprintf(`{"updated": false, "error": %q}`, message))
	} else {
		io.WriteString(w, `{"updated": true}`)
	}
}

func DeleteItemRoute(w http.ResponseWriter, r *http.Request) {
	// Get URL parameter from mux
	vars := mux.Vars(r)
	for k, v := range mux.Vars(r) {
		logrus.Infof("key=%v, value=%v", k, v)
	}
	id := vars["id"]
	err, message := toDoTable.DeleteItem(id)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		io.WriteString(w, fmt.Sprintf(`{"updated": false, "error": %q}`, message))
	} else {
		io.WriteString(w, `{"deleted": true}`)
	}
}

func GetCompletedItems(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Get completed TodoItems")
	completedTodoItems := toDoTable.GetTodoItems(true)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(completedTodoItems)
}

func GetIncompleteItems(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Get Incomplete TodoItems")
	IncompleteTodoItems := toDoTable.GetTodoItems(false)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(IncompleteTodoItems)
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	logrus.Info("API Health is OK")
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}

func Home(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Get index.html")
	p := path.Dir("index.html")
	w.Header().Set("Content-type", "text/html")
	http.ServeFile(w, r, p)
}

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetReportCaller(true)
}

func main() {
	fs := http.FileServer(http.Dir("./resources/"))

	logrus.Info("Starting Todolist API server")
	toDoTable = database.ToDoTable{ToDo: database.GetToDoTable()}
	toDoTable.PrePopulate()
	router := mux.NewRouter()
	router.PathPrefix("/resources/").Handler(http.StripPrefix("/resources/", fs))
	router.HandleFunc("/", Home).Methods("GET")
	router.HandleFunc("/healthz", Healthz).Methods("GET")
	router.HandleFunc("/todo-completed", GetCompletedItems).Methods("GET")
	router.HandleFunc("/todo-incomplete", GetIncompleteItems).Methods("GET")
	router.HandleFunc("/todo", CreateItemRoute).Methods("POST")
	router.HandleFunc("/todo/{id}", UpdateItemRoute).Methods("POST")
	router.HandleFunc("/todo/{id}", DeleteItemRoute).Methods("DELETE")

	handler := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "DELETE", "PATCH", "OPTIONS"},
	}).Handler(router)

	http.ListenAndServe(":8000", handler)
}
