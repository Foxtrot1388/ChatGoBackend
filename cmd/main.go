package main

import (
	"ChatGo/pkg/logging"
	app "ChatGo/server"
)

func main() {
	err := app.Run()
	if err != nil {
		logger := logging.GetLogger()
		logger.Fatal(err)
	}
}
