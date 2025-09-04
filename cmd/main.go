package main

import (
	"mynotes/internal/service"
	"mynotes/pkg/logs"

	"github.com/labstack/echo/v4"
)

func main() {
	logger := logs.NewLogger(false)

	db, err := PostgresConnection()
	if err != nil {
		logger.Fatal(err)
	}

	//defer db.Close()

	svc := service.NewService(db, logger)

	router := echo.New()
	router.Logger = logger

	api := router.Group("/api")

	api.GET("/notes/:id", svc.GetNoteByID)
	api.POST("/notes", svc.CreateNote)
	api.PUT("/notes/:id", svc.UpdateNote)
	api.DELETE("/notes/:id", svc.DeleteNote)
	api.GET("/notes", svc.ListNotes)

	api.POST("/auth/register", svc.Register)
	api.POST("/auth/login", svc.Login)

	router.Logger.Fatal(router.Start(":8000"))
}
