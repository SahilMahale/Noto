package main

import (
	"fmt"
	"os"

	"github.com/SahilMahale/notes-backend/internal/db"
	"github.com/SahilMahale/notes-backend/internal/notes"
	"github.com/SahilMahale/notes-backend/internal/user"
	"github.com/SahilMahale/notes-backend/server"
)

func main() {
	fmt.Println("Connecting to DB")
	db, err := db.NewDbConnection()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to DB")

	ipAddrNPort := os.Getenv("SERVER_BIND_TO")
	if ipAddrNPort == "" {
		ipAddrNPort = "0.0.0.0:8001"
	}
	fmt.Println("Staring server....")
	notesService := server.NewNotesService("Notes app", ipAddrNPort, user.NewUserController(db), notes.NewNoteController(db))
	notesService.StartNotesService()
}
