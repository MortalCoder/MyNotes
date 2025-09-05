package main

import (
	"mynotes/internal/service"
	"mynotes/pkg/logs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	router.Use(middleware.RequestID())

	router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{

		Format: `${time_rfc3339} ` +
			`rid=${id} ` +
			`remote_ip=${remote_ip} ` +
			`method=${method} uri=${uri} ` +
			`status=${status} ` +
			`latency=${latency_human} ` +
			`bytes_in=${bytes_in} bytes_out=${bytes_out} ` +
			`ua="${user_agent}"` + "\n",
	}))

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
