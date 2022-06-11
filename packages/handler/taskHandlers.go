package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type Task struct {
	ID          string `json:"id" db:"ID"`
	Name        string `json:"name" db:"Name"`
	Deadline    string `json:"deadline" db:"Deadline"`
	IsCompleted bool   `json:"isCompleted" db:"IsCompleted"`
	TraQID      string `json:"traQID" db:"TraQID"`
}

func AddTaskHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		task := &Task{}
		err := c.Bind(task)
		traQID := c.Get("UserName").(string)
		if traQID != c.Param("traQID") {
			return echo.NewHTTPError(http.StatusBadRequest, "Do not much traQID")
		}
		if err != nil { // エラーが発生した際
			// fmt.Sprintf("%+v", data): dataをstringに変換
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf(err.Error()))
		}
		taskState := `INSERT INTO tasks(ID, Name, Deadline, IsCompleted, TraQID) VALUES (?,?,?,?,?)`
		_, err2 := db.Exec(taskState, task.ID, task.Name, task.Deadline, task.IsCompleted, traQID)
		if err2 != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf(err2.Error()))
		}

		taskID := task.ID
		//fmt.Println(taskID)
		var rTask Task
		if err := db.Get(&rTask, "SELECT * FROM tasks WHERE ID=?", taskID); errors.Is(err, sql.ErrNoRows) {
			log.Printf("No Such City ID=%s", taskID)
		} else if err != nil {
			log.Fatalf("DB Error: %s", err)
		}

		return c.JSON(http.StatusOK, rTask)
	}
}

func PutTaskHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		task := &Task{}
		err := c.Bind(task)
		traQID := c.Get("UserName").(string)
		if traQID != c.Param("traQID") {
			return echo.NewHTTPError(http.StatusBadRequest, "Do not much traQID")
		}
		if err != nil { // エラーが発生した際
			// fmt.Sprintf("%+v", data): dataをstringに変換
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf(err.Error()))
		}
		taskState := `UPDATE tasks SET Name=?,Deadline=?,IsCompleted=? WHERE ID=? AND TraQID=?`
		_, err2 := db.Exec(taskState, task.Name, task.Deadline, task.IsCompleted, task.ID, traQID)
		if err2 != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf(err2.Error()))
		}

		taskID := task.ID
		//fmt.Println(taskID)
		var rTask Task
		if err := db.Get(&rTask, "SELECT * FROM tasks WHERE ID=?", taskID); errors.Is(err, sql.ErrNoRows) {
			log.Printf("No Such City ID=%s", taskID)
		} else if err != nil {
			log.Fatalf("DB Error: %s", err)
		}

		return c.JSON(http.StatusOK, rTask)
	}
}

func DeleteAllTasksHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		traQID := c.Get("UserName").(string)
		if traQID != c.Param("traQID") {
			return echo.NewHTTPError(http.StatusBadRequest, "Do not much traQID")
		}
		if _, err := db.Exec("DELETE FROM tasks WHERE TraQID=?", traQID); errors.Is(err, sql.ErrNoRows) {
			log.Printf("No Such traQID=%s", traQID)
		} else if err != nil {
			log.Fatalf("DB Error: %s", err)
		}

		return c.String(http.StatusOK, "Delete All Tasks in "+traQID)
	}
}

func DeleteTaskHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		traQID := c.Get("UserName").(string)
		if traQID != c.Param("traQID") {
			return echo.NewHTTPError(http.StatusBadRequest, "Do not much traQID")
		}
		ID := c.Param("ID")
		//fmt.Println(cityID)
		if _, err := db.Exec("DELETE FROM tasks WHERE ID=? AND TraQID=?", ID, traQID); errors.Is(err, sql.ErrNoRows) {
			log.Printf("No Such task ID=%s in this user", ID)
		} else if err != nil {
			log.Fatalf("DB Error: %s", err)
		}

		return c.String(http.StatusOK, "Delete task ID:"+ID)
	}
}

func GetTaskHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		tasks := []Task{}
		traQID := c.Get("UserName").(string)
		if traQID != c.Param("traQID") {
			return echo.NewHTTPError(http.StatusBadRequest, "Do not much traQID")
		}
		if err := db.Select(&tasks, "SELECT * FROM tasks WHERE traQID=?", traQID); errors.Is(err, sql.ErrNoRows) {
			log.Printf("No Such traQID=%s", traQID)
		} else if err != nil {
			log.Fatalf("DB Error: %s", err)
		}
		return c.JSON(http.StatusOK, tasks)
	}
}
