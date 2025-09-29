package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func NoteRoutes(v1 fiber.Router, noteService service.NoteService, userService service.UserService) {
	noteController := controller.NewNoteController(noteService)
	notes := v1.Group("/notes")

	notes.Use(m.Auth(userService))

	notes.Post("/", noteController.CreateNote)
	notes.Get("/", noteController.GetNotes)
	notes.Get("/:noteId", noteController.GetNote)
	notes.Patch("/:noteId", noteController.UpdateNote)
	notes.Delete("/:noteId", noteController.DeleteNote)
	notes.Post("/:noteId/upload", noteController.UploadFile)
}
