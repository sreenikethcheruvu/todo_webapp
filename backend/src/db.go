package src

import (
	"encoding/json"
	"github.com/dgraph-io/badger/v4"
	"log"
)

var db *badger.DB

func InitDB() {
	opts := badger.DefaultOptions("./data").WithLogger(nil)
	var err error
	db, err = badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
}

func SaveTodo(todo Todo) error {
	return db.Update(func(txn *badger.Txn) error {
		data, err := json.Marshal(todo)
		if err != nil {
			return err
		}
		return txn.Set([]byte(todo.ID), data)
	})
}

func GetTodo(id string) (Todo, error){
	var todo Todo
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(id))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &todo)
		})
	})
	return todo, err
}

func GetAllTodos() ([]Todo, error) {
	var todos []Todo

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = true
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var todo Todo
				if err := json.Unmarshal(val, &todo); err != nil {
					return err
				}
				todos = append(todos, todo)
				return nil
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	return todos, err
}


func UpdateTodoStatus(id string, completed bool) error {
	todo, err := GetTodo(id)
	if err != nil {
		return err 
	}

	todo.Completed = completed

	return SaveTodo(todo)
}

func RenameTodo(id string, newTitle string) error {
	todo, err := GetTodo(id)
	if err != nil {
		return err
	}

	todo.Title = newTitle
	return SaveTodo(todo)
}

func DeleteTodo(id string) error {
	return db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(id))
	})
}
