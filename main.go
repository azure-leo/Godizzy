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

func getTaskHandler(c echo.Context) error {
	var tasks []Task
	if err := db.Find(&tasks).Error; err != nil {
		return c.String(http.StatusInternalServerError, "Could not get tasks")
	}
	return c.JSON(http.StatusOK, tasks)
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

	e.Logger.Fatal(e.Start(":9090"))
}
