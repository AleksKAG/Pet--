package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Task — модель задачи
type Task struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Task      string         `json:"task"`
	IsDone    bool           `json:"is_done"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}


var db *gorm.DB

func main() {
	var err error

	// Подключение к PostgreSQL
	dsn := "host=localhost user=postgres password=123456 dbname=tasks_db port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // Отключить множественные имена таблиц
		},
	})
	if err != nil {
		panic("не удалось подключиться к PostgreSQL")
	}

	// Миграция таблицы
	db.AutoMigrate(&Task{})

	// Echo router
	e := echo.New()

	e.POST("/tasks", createTask)
	e.GET("/tasks", getAllTasks)
	e.PATCH("/tasks/:id", updateTask)
	e.DELETE("/tasks/:id", deleteTask)

	e.Logger.Fatal(e.Start(":8080"))
}

// Создание новой задачи
func createTask(c echo.Context) error {
	var task Task
	if err := c.Bind(&task); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Некорректный JSON"})
	}
	db.Create(&task)
	return c.JSON(http.StatusCreated, task)
}

// Получение всех задач
func getAllTasks(c echo.Context) error {
	var tasks []Task
	if err := db.Where("deleted_at IS NULL").Find(&tasks).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Ошибка при получении задач"})
	}
	return c.JSON(http.StatusOK, tasks)
}


// Обновление задачи по ID
func updateTask(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var task Task

	if err := db.First(&task, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Задача не найдена"})
	}

	var input Task
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Ошибка JSON"})
	}

	task.Task = input.Task
	task.IsDone = input.IsDone
	db.Save(&task)

	return c.JSON(http.StatusOK, task)
}

// Удаление задачи (мягкое)
func deleteTask(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var task Task

	if err := db.First(&task, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Задача не найдена"})
	}

	db.Delete(&task)
	return c.NoContent(http.StatusNoContent)
}
