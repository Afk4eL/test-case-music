package main

import (
	"os"
	"os/signal"
	"syscall"

	"test-case/internal/app"
)

func main() {
	application := app.New()

	application.SetConfig()

	go application.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.Stop()
}
