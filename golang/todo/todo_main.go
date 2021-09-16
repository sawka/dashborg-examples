package main

import (
	"fmt"
	"reflect"

	"github.com/sawka/dashborg-go-sdk/pkg/dash"
)

type TodoItem struct {
	Id       int
	TodoType string
	Item     string
	Done     bool
}

type ServerTodoModel struct {
	TodoList []*TodoItem
	NextId   int
}

type TodoAppState struct {
	TodoType string `json:"todotype"`
	NewTodo  string `json:"newtodo"`
}

func (m *ServerTodoModel) AddTodo(req *dash.AppRequest, state *TodoAppState) error {
	if state.NewTodo == "" {
		return fmt.Errorf("Please enter a Todo Item")
	}
	if state.TodoType == "" {
		return fmt.Errorf("Please select a Todo Type")
	}
	m.TodoList = append(m.TodoList, &TodoItem{Id: m.NextId, Item: state.NewTodo, TodoType: state.TodoType})
	m.NextId++
	req.InvalidateData("/@app:get-todo-list")
	return nil
}

func (m *ServerTodoModel) MarkTodoDone(req *dash.AppRequest, todoId int) error {
	for _, todoItem := range m.TodoList {
		if todoItem.Id == todoId {
			todoItem.Done = true
		}
	}
	req.InvalidateData("/@app:get-todo-list")
	return nil
}

func (m *ServerTodoModel) RemoveTodo(req *dash.AppRequest, todoId int) error {
	newList := make([]*TodoItem, 0)
	for _, todoItem := range m.TodoList {
		if todoItem.Id == todoId {
			continue
		}
		newList = append(newList, todoItem)
	}
	m.TodoList = newList
	req.InvalidateData("/@app:get-todo-list")
	return nil
}

func (m *ServerTodoModel) GetTodoList() (interface{}, error) {
	return m.TodoList, nil
}

func main() {
	cfg := &dash.Config{AnonAcc: true, AutoKeygen: true}
	client, err := dash.ConnectClient(cfg)
	if err != nil {
		panic(err)
	}
	tm := &ServerTodoModel{NextId: 1}
	app := client.AppClient().NewApp("todo")
	app.WatchHtmlFile("panels/todo.html", nil)
	app.Runtime().SetAppStateType(reflect.TypeOf(&TodoAppState{}))
	app.Runtime().Handler("add-todo", tm.AddTodo)
	app.Runtime().Handler("mark-todo-done", tm.MarkTodoDone)
	app.Runtime().Handler("remove-todo", tm.RemoveTodo)
	app.Runtime().PureHandler("get-todo-list", tm.GetTodoList)
	err = client.AppClient().WriteAndConnectApp(app)
	if err != nil {
		panic(err)
	}
	client.WaitForShutdown()
}
