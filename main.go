package main

import (
	"log"
	"os"

	"gitlab.com/sckacr/calltaxi/app"
	"gitlab.com/sckacr/calltaxi/config"
	"gitlab.com/sckacr/calltaxi/router"
)

func main() {
	if err := app.Init(config.New()); err != nil {
		log.Fatalf("%+v\n", err)
	}

	log.Fatal(router.New().Run(":" + os.Getenv("PORT")))
}
