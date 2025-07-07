package src

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

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

	if err := c.ShouldBindJSON(&input); err != nil || input.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	todo := Todo{
		ID:        uuid.New().String(),
		Title:     input.Title,
		Completed: false,
	}

	if err := SaveTodo(todo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

func getTodoHandler(c *gin.Context) {
	id := c.Param("id")
	todo, err := GetTodo(id)
	if err == ErrTodoNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func deleteTodoHandler(c *gin.Context) {
	id := c.Param("id")
	err := DeleteTodo(id)
	if err == ErrTodoNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted"})
}

func renameTodoHandler(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Title string `json:"title"`
	}

	if err := c.ShouldBindJSON(&input); err != nil || input.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	err := RenameTodo(id, input.Title)
	if err == ErrTodoNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo title updated"})
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

	err := UpdateTodoStatus(id, input.Completed)
	if err == ErrTodoNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo status updated"})
}