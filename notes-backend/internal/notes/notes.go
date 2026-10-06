package notes

import (
	"github.com/SahilMahale/notes-backend/internal/db"
	"github.com/SahilMahale/notes-backend/internal/helper"
	"github.com/SahilMahale/notes-backend/models"
	"github.com/google/uuid"
)

type NotesController struct {
	DbInterface db.DbConnection
}

type NotesOps interface {
	CreateNote(title, body string) (string, helper.MyHTTPErrors)
	UpdateNote(string, models.NotePatchRequest) (models.NotePatchResp, helper.MyHTTPErrors)
	DeleteNote(noteId string) helper.MyHTTPErrors
	GetAllNotes() (models.NotesResp, helper.MyHTTPErrors)
	GetNote(noteID string) (models.NoteResp, helper.MyHTTPErrors)
}

func NewNoteController(db db.DbConnection) NotesController {
	return NotesController{
		DbInterface: db,
	}
}

func (b NotesController) CreateNote(title, body string) (string, helper.MyHTTPErrors) {
	noteId := uuid.NewString()
	note := db.Note{
		ID: noteId,
		Title:  title,
		Body:   body,
	}
	if err := b.DbInterface.Db.Create(&note); err.Error != nil {
		return "", helper.ErrorMatch(err.Error)
	}
	return noteId, helper.MyHTTPErrors{
		Err: nil,
	}
}

func (n NotesController) GetAllNotes() (models.NotesResp, helper.MyHTTPErrors) {
	notesList := models.NotesResp{}
	notes := []db.Note{}
	res := n.DbInterface.Db.Find(&notes)
	if res.Error != nil {
		myerr := helper.ErrorMatch(res.Error)
		return nil, myerr
	}
	for _, note := range notes {
		notesList = append(notesList, models.NoteResp{
			NoteID: note.ID,
			Title:  note.Title,
			Body:   note.Body,
		})
	}
	return notesList, helper.MyHTTPErrors{
		Err: nil,
	}
}

func (n NotesController) DeleteNote(noteId string) helper.MyHTTPErrors {
	note := db.Note{
		ID: noteId,
	}
	err := n.DbInterface.Db.First(&note)
	if err.Error != nil {
		return helper.ErrorMatch(err.Error)
	}

	err = n.DbInterface.Db.Delete(&note)
	if err.Error != nil {
		return helper.ErrorMatch(err.Error)
	}
	return helper.MyHTTPErrors{
		Err: nil,
	}
}

// UpdateNote method to update individual fields of a note
func (n NotesController) UpdateNote(noteID string, note models.NotePatchRequest) (models.NotePatchResp, helper.MyHTTPErrors) {
	res := n.DbInterface.Db.Model(&db.Note{}).Where("note_id = ?", noteID).Updates(&note)
	if res.Error != nil {
		return models.NotePatchResp{}, helper.ErrorMatch(res.Error)
	}
	return models.NotePatchResp{
			NoteID: noteID,
			Title:  note.Title,
			Body:   note.Body,
		}, helper.MyHTTPErrors{
			Err: nil,
		}
}

func (n NotesController) GetNote(noteID string) (models.NoteResp, helper.MyHTTPErrors) {
	note := db.Note{
		ID: noteID,
	}
	err := n.DbInterface.Db.First(&note)
	if err.Error != nil {
		return models.NoteResp{}, helper.ErrorMatch(err.Error)
	}
	return models.NoteResp{
			NoteID: note.ID,
			Title:  note.Title,
			Body:   note.Body,
		}, helper.MyHTTPErrors{
			Err: nil,
		}
}
