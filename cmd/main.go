package main

import (
	"mynotes/internal/service"
	"mynotes/pkg/logs"

	"github.com/labstack/echo/v4"
)

func main() {
	// 1) логгер (как раньше)
	logger := logs.NewLogger(false)

	// 2) подключение к БД (как у тебя — одной функцией)
	db, err := PostgresConnection()
	if err != nil {
		logger.Fatal(err)
	}

	// 3) сервис-слой
	svc := service.NewService(db, logger)

	// 4) роутер и группа /api
	router := echo.New()
	api := router.Group("/api")

	//api.GET("/healthz", svc.Health)

	//api.GET("/notes/:id", svc.GetNoteByID)
	//api.POST("/notes", svc.CreateNote)
	//api.PUT("/notes/:id", svc.UpdateNote)
	//api.DELETE("/notes/:id", svc.DeleteNote)
	//api.GET("/notes", svc.ListNotes) // ?page, ?limit

	// api.POST("/auth/register", svc.Register)
	// api.POST("/auth/login", svc.Login)

	// 8) стартуем на 8000, как раньше
	router.Logger.Fatal(router.Start(":8000"))
}
