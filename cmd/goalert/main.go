package main

import (
	"os"
	_ "time/tzdata"

	"github.com/breathbath/goalert/app"
	"github.com/breathbath/goalert/util/log"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	l := log.NewLogger()
	ctx := l.BackgroundContext()
	err := app.RootCmd.ExecuteContext(ctx)
	if err != nil {
		log.Log(ctx, err)
		os.Exit(1)
	}
}
