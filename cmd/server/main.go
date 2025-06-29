package main

import (
	"itk-assignment/server"
	"log"
)

func main() {
	if err := server.Start(); err != nil {
		log.Fatalln(err)
	}
}
