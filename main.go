package main

import (
    "github.com/labstack/echo/v4"
    "net/http"
    "strconv"
)

type Task struct {
    ID     int    `json:"id"`
    Name   string `json:"task"`
}

var tasks = make(map[int]Task) // хранилище задач
var nextID = 1                 // для генерации уникальных ID

func main() {
//Создаем новое приложение Echo
    e := echo.New()

    // POST /task — создаёт задачу
    e.POST("/task", func(c echo.Context) error {
        var newTask Task
        if err := c.Bind(&newTask); err != nil {
            return err
        }

        newTask.ID = nextID
        tasks[nextID] = newTask
        nextID++

        return c.JSON(http.StatusCreated, newTask)
    })

    // GET / — получить все задачи
    e.GET("/", func(c echo.Context) error {
        var result []Task
        for _, task := range tasks {
            result = append(result, task)
        }
        return c.JSON(http.StatusOK, result)
    })

    // PATCH /task/:id — обновить задачу по ID
    e.PATCH("/task/:id", func(c echo.Context) error {
        id, _ := strconv.Atoi(c.Param("id"))
        var updated Task
        if err := c.Bind(&updated); err != nil {
            return err
        }

        task, exists := tasks[id]
        if !exists {
            return c.String(http.StatusNotFound, "Задача не найдена")
        }

        task.Name = updated.Name
        tasks[id] = task

        return c.JSON(http.StatusOK, task)
    })

    // DELETE /task/:id — удалить задачу по ID
    e.DELETE("/task/:id", func(c echo.Context) error {
        id, _ := strconv.Atoi(c.Param("id"))
        if _, exists := tasks[id]; !exists {
            return c.String(http.StatusNotFound, "Задача не найдена")
        }
        delete(tasks, id)
        return c.NoContent(http.StatusNoContent)
    })

    e.Start(":8080")
}
