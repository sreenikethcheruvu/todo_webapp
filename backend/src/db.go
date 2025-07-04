package src

import (
    "context"
    "errors"
    "sync"

    "github.com/google/uuid"
    pb "todoapp-backend/src/pb"
)

type todoStore struct {
    sync.RWMutex
    todos map[string]*pb.Todo
}

var store = todoStore{
    todos: make(map[string]*pb.Todo),
}

type TodoServiceServer struct {
    pb.UnimplementedTodoServiceServer
}

func (s *TodoServiceServer) CreateTodo(ctx context.Context, req *pb.CreateTodoRequest) (*pb.CreateTodoResponse, error) {
    id := uuid.New().String()
    todo := &pb.Todo{
        Id:        id,
        Title:     req.Title,
        Completed: false,
    }

    store.Lock()
    store.todos[id] = todo
    store.Unlock()

    return &pb.CreateTodoResponse{Todo: todo}, nil
}

func (s *TodoServiceServer) GetAllTodos(ctx context.Context, _ *pb.Empty) (*pb.GetAllTodosResponse, error) {
    store.RLock()
    defer store.RUnlock()

    var todos []*pb.Todo
    for _, todo := range store.todos {
        todos = append(todos, todo)
    }
    return &pb.GetAllTodosResponse{Todos: todos}, nil
}

func (s *TodoServiceServer) GetTodo(ctx context.Context, req *pb.GetTodoRequest) (*pb.GetTodoResponse, error) {
    store.RLock()
    defer store.RUnlock()

    todo, ok := store.todos[req.Id]
    if !ok {
        return nil, errors.New("todo not found")
    }
    return &pb.GetTodoResponse{Todo: todo}, nil
}

func (s *TodoServiceServer) RenameTodo(ctx context.Context, req *pb.RenameTodoRequest) (*pb.Empty, error) {
    store.Lock()
    defer store.Unlock()

    todo, ok := store.todos[req.Id]
    if !ok {
        return nil, errors.New("todo not found")
    }
    todo.Title = req.Title
    return &pb.Empty{}, nil
}

func (s *TodoServiceServer) UpdateTodoStatus(ctx context.Context, req *pb.UpdateTodoStatusRequest) (*pb.Empty, error) {
    store.Lock()
    defer store.Unlock()

    todo, ok := store.todos[req.Id]
    if !ok {
        return nil, errors.New("todo not found")
    }
    todo.Completed = req.Completed
    return &pb.Empty{}, nil
}

func (s *TodoServiceServer) DeleteTodo(ctx context.Context, req *pb.DeleteTodoRequest) (*pb.Empty, error) {
    store.Lock()
    defer store.Unlock()

    if _, ok := store.todos[req.Id]; !ok {
        return nil, errors.New("todo not found")
    }
    delete(store.todos, req.Id)
    return &pb.Empty{}, nil
}
