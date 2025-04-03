package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

var task string

type requestBody struct {
	Task string `json:task`
}

func postTaskHandler(c echo.Context) error {
	var request requestBody
	if err := c.Bind(&request); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request")
	}
	task = request.Task
	return c.String(http.StatusOK, "Task received")
}

func getTaskHandler(c echo.Context) error {
	return c.String(http.StatusOK, "hello, "+task)
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}

func main() {
	e := echo.New()
	e.POST("/task", postTaskHandler)
	e.GET("/task", getTaskHandler)

	e.Logger.Fatal(e.Start(":9090"))
}
