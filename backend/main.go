package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"backend/models"
	"backend/server"
)

func init_log() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Println("Failed to open log file:", err)
	}
	log.SetOutput(logFile)
}

func main() {
	init_log()
	service, err := models.New()
	if err != nil {
		log.Fatalf("Failed to initialize model for operating all service, %s\n", err)
	}
	server := server.NewServer(service)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to listen for http server, %s\n", err)
	}
	fmt.Printf("%v", server.Addr)
	log.Print("start serving")
}
