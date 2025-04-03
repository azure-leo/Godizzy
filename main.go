package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

var db *gorm.DB

func initDB() {
	dsn := "host=localhost user=postgres password=Felixculpa.16 dbname=postgres port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to database:  %v", err)
	}
}

type requestBody struct {
	Task   string `json:"task""`
	IsDone bool   `json:"is_done"`
}

func postTaskHandler(c echo.Context) error {
	var request requestBody
	if err := c.Bind(&request); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	newTask := Task{
		Task:   request.Task,
		IsDone: request.IsDone,
	}

	if err := db.Create(&newTask).Error; err != nil {
		return c.String(http.StatusInternalServerError, "Could not create task")
	}

	return c.JSON(http.StatusOK, newTask)
}

func updateTaskHandler(c echo.Context) error {
	var task Task
	id := c.Param("id")

	var request requestBody
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid patch request")
	}

	if err := db.First(&task, id).Error; err != nil {
		return c.String(http.StatusNotFound, "Could not find task")
	}
	task.Task = request.Task
	task.IsDone = request.IsDone

	if err := db.Save(&task).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, "Could not update task")
	}

	return c.JSON(http.StatusOK, task)
}

func getTaskHandler(c echo.Context) error {
	var tasks []Task
	if err := db.Find(&tasks).Error; err != nil {
		return c.String(http.StatusInternalServerError, "Could not get tasks")
	}
	return c.JSON(http.StatusOK, tasks)
}

func deleteTaskHandler(c echo.Context) error {
	id := c.Param("id")
	if err := db.Delete(&Task{}, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, "Could not delete task")
	}
	return c.JSON(http.StatusOK, "Task deleted")
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}

func main() {
	initDB()
	db.AutoMigrate(&Task{})
	e := echo.New()
	e.POST("/task", postTaskHandler)
	e.GET("/task", getTaskHandler)
	e.PATCH("/task/:id", updateTaskHandler)
	e.DELETE("/task/:id", deleteTaskHandler)

	e.Logger.Fatal(e.Start(":9090"))
}
