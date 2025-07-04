package src

import (
    "context"
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "google.golang.org/grpc"
    pb "todoapp-backend/src/pb"
)

var grpcClient pb.TodoServiceClient

func SetupRoutes() *gin.Engine {
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to connect to gRPC server: %v", err)
    }
    grpcClient = pb.NewTodoServiceClient(conn)

    r := gin.Default()

    r.GET("/todos", getAllTodosHandler)
    r.GET("/todos/:id", getTodoHandler)
    r.POST("/todos", createTodoHandler)
    r.PUT("/todos/:id/name", renameTodoHandler)
    r.PUT("/todos/:id/status", updateTodoStatusHandler)
    r.DELETE("/todos/:id", deleteTodoHandler)

    return r
}

func createTodoHandler(c *gin.Context) {
    var input struct {
        Title string `json:"title"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    resp, err := grpcClient.CreateTodo(context.Background(), &pb.CreateTodoRequest{
        Title: input.Title,
    })
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, resp.Todo)
}

func getAllTodosHandler(c *gin.Context) {
    resp, err := grpcClient.GetAllTodos(context.Background(), &pb.Empty{})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, resp.Todos)
}

func getTodoHandler(c *gin.Context) {
    id := c.Param("id")
    resp, err := grpcClient.GetTodo(context.Background(), &pb.GetTodoRequest{Id: id})
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
        return
    }
    c.JSON(http.StatusOK, resp.Todo)
}

func renameTodoHandler(c *gin.Context) {
    id := c.Param("id")
    var input struct {
        Title string `json:"title"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    _, err := grpcClient.RenameTodo(context.Background(), &pb.RenameTodoRequest{Id: id, Title: input.Title})
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Todo renamed"})
}

func updateTodoStatusHandler(c *gin.Context) {
    id := c.Param("id")
    var input struct {
        Completed bool `json:"completed"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    _, err := grpcClient.UpdateTodoStatus(context.Background(), &pb.UpdateTodoStatusRequest{Id: id, Completed: input.Completed})
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Todo status updated"})
}

func deleteTodoHandler(c *gin.Context) {
    id := c.Param("id")
    _, err := grpcClient.DeleteTodo(context.Background(), &pb.DeleteTodoRequest{Id: id})
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Todo deleted"})
}
