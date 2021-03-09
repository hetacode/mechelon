package main

import (
	"log"
	"os"

	svvcontainer "github.com/hetacode/mechelon/services/services-view/container"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../../.env.dev")

	log.Println("services view svc is starting")
	waitCh := make(<-chan os.Signal)

	_ = svvcontainer.NewContainer()

	log.Println("services view is running")
	<-waitCh
}
