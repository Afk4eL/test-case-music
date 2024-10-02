package router

import (
	"test-case/internal/server/handlers"
	middleware_logger "test-case/internal/server/middlewares/logger"
	"test-case/storage/repos"

	_ "test-case/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Song Library API
// @version 1.0
// @description This is an API for managing songs in a library.
// @host localhost:8080
// @BasePath /

func SetupRouter(repos repos.SongRepository) *gin.Engine {
	router := gin.Default()

	handler := handlers.NewSongHandler(repos)

	router.Use(middleware_logger.RequestLogger())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/get-songs", handler.GetSongs)
	router.GET("/get-song-text", handler.GetSongText)
	router.DELETE("/delete-song", handler.DeleteSong)
	router.POST("/update-song", handler.UpdateSong)
	router.POST("/add-song", handler.AddSong)

	//ДЛЯ ДЕБАГА
	router.GET("/info", func(c *gin.Context) {
		c.JSON(200, gin.H{"releaseDate": "16.07.2006", "text": "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight", "link": "https://www.youtube.com/watch?v=Xsp3_a-PMTw"})
	})

	return router
}
