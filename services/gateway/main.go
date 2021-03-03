package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../../.env.dev")

	log.Println("api gateway svc is starting")
	waitCh := make(<-chan os.Signal)

	log.Println("api gateway svc is running")
	<-waitCh
}
