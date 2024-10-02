package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"test-case/internal/config"
	"test-case/internal/server/router"
	"test-case/internal/utils/logger"
	"test-case/storage/postgres"
	"test-case/storage/repos"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

type App struct {
	Cfg      config.Config
	Storage  *postgres.Database
	SongRepo repos.SongRepository
	Router   *gin.Engine
	Server   *http.Server
}

func (app *App) readConfig() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("Usage go run <path to main.go> [arguments] \n Required arguments: \n - Path to config file")
		os.Exit(1)
	}

	app.Cfg = config.ReadConfig(args[0])
}

func (app *App) SetConfig() {
	app.readConfig()

	storage, err := postgres.New(app.Cfg)
	if err != nil {

		os.Exit(1)
	}
	app.Storage = storage

	app.SongRepo = repos.NewSongRepository(app.Storage.Database)

	app.Router = router.SetupRouter(app.SongRepo)

	app.Server = &http.Server{
		Addr:    app.Cfg.Address,
		Handler: app.Router,
	}
}

func (app *App) Run() {
	const op = "app.Run"

	logger.Logger.Info().Str("address:", app.Server.Addr).Msg("Server started")
	if err := app.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Logger.Fatal().Msg(fmt.Sprint("Fatal error", op, err.Error()))
		return
	}
}

func (app *App) Stop() {
	const op = "app.Stop"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	app.Storage.Stop()

	if err := app.Server.Shutdown(ctx); err != nil {
		logger.Logger.Warn().Msg(fmt.Sprint("Server forced to shutdown", op, err.Error()))
		return
	}

	logger.Logger.Info().Msg("Server stopped")
}

func New() *App {
	return &App{}
}
