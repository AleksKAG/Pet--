package main

import (
    "net/http"
    "github.com/labstack/echo/v4"
)

// Глобальная переменная для хранения task
var task string

// Структура для чтения JSON из POST-запроса
type RequestBodyTask struct {
    Task string `json:"task"`
}

func main() {
     e := echo.New() // Создаем новый экземпляр Echo



    // Регистрируем маршруты
    e.POST("/task", handleTask)  // POST /task для получения нового task


    e.GET("/", handleRoot)       // GET / возвращает "hello {task}"

    // Запускаем сервер на порту 8080
    e.Logger.Fatal(e.Start(":8080"))
}

// Обработчик POST-запроса на /task
func handleTask(c echo.Context) error {
    // Создаем переменную структуры, в которую будет декодирован JSON
    var req RequestBodyTask

    // Парсим JSON из тела запроса в переменную req
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Некорректный JSON",
        })
    }

    // Сохраняем значение task в глобальную переменную
    task = req.Task

    // Отправляем ответ клиенту
    return c.JSON(http.StatusOK, map[string]string{
        "message": "Task сохранён",
        "task":    task,
    })
}

// Обработчик GET-запроса на /
func handleRoot(c echo.Context) error {
    // Возвращаем сохранённый task
    return c.String(http.StatusOK, "hello " + task + "!")
}
