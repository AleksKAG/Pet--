package main

import (
	"net/http"          // Для HTTP-ответов (статусы и методы)
	"strconv"           // Для преобразования строк в числа (например, id из URL)

	"github.com/labstack/echo/v4"     // Web-фреймворк Echo
	"gorm.io/driver/sqlite"           // SQLite-драйвер для GORM
	"gorm.io/gorm"                    // GORM — ORM-библиотека для работы с базой данных
)

// Структура задачи. Она станет таблицей в базе данных.
type Task struct {
	ID      uint           `json:"id" gorm:"primaryKey"` // Уникальный ID (Primary Key)
	Task    string         `json:"task"`                 // Название задачи
	IsDone  bool           `json:"is_done"`              // Статус выполнения
	Deleted gorm.DeletedAt `gorm:"index"`                // Для мягкого удаления (soft delete)
}

// Глобальная переменная для подключения к БД
var db *gorm.DB

func main() {
	var err error

	// Подключаемся к базе данных SQLite (файл будет создан, если его нет)
	db, err = gorm.Open(sqlite.Open("tasks.db"), &gorm.Config{})
	if err != nil {
		panic("не удалось подключиться к базе")
	}

	// Создаем таблицу в БД по структуре Task, если её ещё нет
	db.AutoMigrate(&Task{})

	// Создаем сервер Echo
	e := echo.New()

	// Регистрируем маршруты (ручки)

	// POST /task — создать новую задачу
	e.POST("/task", createTask)

	// GET /tasks — получить список всех задач
	e.GET("/tasks", getAllTasks)

	// PATCH /task/:id — обновить задачу по ID
	e.PATCH("/task/:id", updateTask)

	// DELETE /task/:id — мягко удалить задачу по ID
	e.DELETE("/task/:id", deleteTask)

	// Запускаем сервер на порту 8080
	e.Logger.Fatal(e.Start(":8080"))
}

// Обработчик POST /task — создаёт новую задачу
func createTask(c echo.Context) error {
	var task Task

	// Читаем JSON из тела запроса в структуру
	if err := c.Bind(&task); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Некорректный JSON"})
	}

	// Сохраняем задачу в базу данных
	db.Create(&task)

	// Отправляем обратно созданную задачу
	return c.JSON(http.StatusCreated, task)
}

// Обработчик GET /tasks — возвращает все задачи
func getAllTasks(c echo.Context) error {
	var tasks []Task

	// Читаем все задачи из базы (кроме удалённых)
	db.Find(&tasks)

	// Возвращаем JSON со списком задач
	return c.JSON(http.StatusOK, tasks)
}

// Обработчик PATCH /task/:id — обновляет задачу
func updateTask(c echo.Context) error {
	// Получаем ID из параметра URL и преобразуем его в число
	id, _ := strconv.Atoi(c.Param("id"))
	var task Task

	// Ищем задачу по ID
	if err := db.First(&task, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Задача не найдена"})
	}

	// Читаем обновлённые данные из тела запроса
	var input Task
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Ошибка JSON"})
	}

	// Обновляем поля
	task.Task = input.Task
	task.IsDone = input.IsDone
	db.Save(&task)

	// Возвращаем обновлённую задачу
	return c.JSON(http.StatusOK, task)
}

// Обработчик DELETE /task/:id — мягко удаляет задачу
func deleteTask(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var task Task

	// Ищем задачу
	if err := db.First(&task, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Задача не найдена"})
	}

	// Мягкое удаление (GORM проставит флаг DeletedAt)
	db.Delete(&task)

	// Возвращаем пустой ответ (204 No Content)
	return c.NoContent(http.StatusNoContent)
}

