package server

import (
	"crypto/rsa"
	"fmt"

	"github.com/SahilMahale/notes-backend/internal/notes"
	"github.com/SahilMahale/notes-backend/internal/user"
	"github.com/SahilMahale/notes-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

type MyCustomClaims struct {
	Name string `json:"name"`
	jwt.RegisteredClaims
}

type notesService struct {
	app       *fiber.App
	notesCtrl notes.NotesOps
	userCtrl  user.UserOps
	ip        string
}

type notesServicer interface {
	GetNotes(c *fiber.Ctx) error
	DeleteNote(c *fiber.Ctx) error
	CreateNote(c *fiber.Ctx) error
	UpdateNote(c *fiber.Ctx) error
	StartNotesService(c *fiber.Ctx) error
}

func NewNotesService(appname, ip string, userCtrl user.UserOps, notesCtrl notes.NotesOps) notesService {
	return notesService{
		app: fiber.New(fiber.Config{
			AppName:       appname,
			StrictRouting: true,
			ServerHeader:  "NotesService",
		}),
		ip:        ip,
		userCtrl:  userCtrl,
		notesCtrl: notesCtrl,
	}
}

func (B *notesService) GetNotes(c *fiber.Ctx) error {
	notesList, err := B.notesCtrl.GetAllNotes()
	if err.Err != nil {
		return c.Status(err.HttpCode).SendString(err.Err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(notesList)
}

func (B *notesService) GetNoteByID(c *fiber.Ctx) error {
	noteID := ""
	if noteID = c.Params("noteID"); noteID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Error: notedID not specified")
	}
	note, err := B.notesCtrl.GetNote(noteID)
	if err.Err != nil {
		return c.Status(err.HttpCode).SendString(err.Err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(note)
}

func (B *notesService) DeleteNote(c *fiber.Ctx) error {
	noteID := ""
	if noteID = c.Params("noteID"); noteID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Error: notedID not specified")
	}
	err := B.notesCtrl.DeleteNote(noteID)
	if err.Err != nil {
		return c.Status(err.HttpCode).SendString(err.Err.Error())
	}
	return c.Status(fiber.StatusOK).SendString(fmt.Sprintf("Note with notedID: %s is successfully deleted.\n", noteID))
}

func (B *notesService) CreateNote(c *fiber.Ctx) error {
	note := new(models.NoteRequest)

	if err := c.BodyParser(note); err != nil {
		return err
	}
	noteID, err := B.notesCtrl.CreateNote(note.Title, note.Body)
	if err.Err != nil {
		return c.Status(err.HttpCode).SendString(err.Err.Error())
	}

	noteResp := models.NoteResp{
		NoteID: noteID,
		Title:  note.Title,
		Body:   note.Body,
	}
	return c.Status(fiber.StatusAccepted).JSON(noteResp)
}

func (B *notesService) UpdateNote(c *fiber.Ctx) error {
	noteID := ""
	if noteID = c.Params("noteID"); noteID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Error: notedID not specified")
	}

	note := new(models.NotePatchRequest)

	if err := c.BodyParser(note); err != nil {
		return err
	}
	noteResp, err := B.notesCtrl.UpdateNote(noteID, *note)
	if err.Err != nil {
		return c.Status(err.HttpCode).SendString(err.Err.Error())
	}
	return c.Status(fiber.StatusAccepted).JSON(noteResp)
}

func (B *notesService) CreateUser(c *fiber.Ctx) error {
	u := new(models.UserSignup)

	if err := c.BodyParser(u); err != nil {
		return err
	}

	errP := B.userCtrl.CreateUser(u.Username, u.Email, u.Password)
	if errP.Err != nil {
		return c.Status(errP.HttpCode).SendString(errP.Err.Error())
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (B *notesService) LoginUser(c *fiber.Ctx) error {
	u := new(models.UserSignin)
	if err := c.BodyParser(u); err != nil {
		return err
	}

	userID, err := B.userCtrl.LoginUser(u.Username, u.Password)
	if err.Err != nil {
		return c.Status(err.HttpCode).SendString(err.Err.Error())
	}

	// Create a token based on user
	atoken, errp := makeTokenWithClaims(userID,u.Username)

	if errp != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(errp.Error())
	}
	
	resp := models.JwtResp{
		Authtoken: atoken,
	}
	return c.Status(fiber.StatusAccepted).JSON(resp)
}

// Stubbed version of StartNotesService is tested
//
//nolint:coverage
func (B *notesService) StartNotesService() {
	B.initMiddleware()
	// Unauthenticated routes
	B.initAuth()
	userGroup := B.app.Group("/user")
	userGroup.Post("/signup", B.CreateUser)
	userGroup.Post("/signin", B.LoginUser)

	B.app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Booking APP Service is Running!")
	})
	// authenticated routes

	notesGroup := B.app.Group("/notes")
	notesGroup.Post("/create", B.CreateNote)
	notesGroup.Patch("/update/:noteID", B.UpdateNote)
	notesGroup.Delete("/delete/:noteID", B.DeleteNote)
	notesGroup.Get("/get", B.GetNotes)
	notesGroup.Get("/get/:noteId", B.GetNoteByID)

	err := B.app.Listen(B.ip)
	if err != nil {
		log.Error(err)
		return
	}
}
