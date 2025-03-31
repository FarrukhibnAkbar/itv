package main

import (
	"itv/monorepo/movie_service/app"
)

func main() {
	app := app.New()
	app.Run()
}
