package controller

import (
	"app/src/model"
	"app/src/service"
	"app/src/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type NoteController struct {
	NoteService service.NoteService
}

func NewNoteController(noteService service.NoteService) *NoteController {
	return &NoteController{NoteService: noteService}
}

// @Tags         Notes
// @Summary      Create a new note
// @Description  Creates a new note associated with the authenticated user.
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      validation.CreateNote  true  "Request body"
// @Success      201      {object}  example.CreateNoteResponse
// @Failure      400      {object}  example.Unauthorized  "Bad Request (validation error)"
// @Failure      401      {object}  example.Unauthorized
// @Router       /notes [post]
func (nc *NoteController) CreateNote(c *fiber.Ctx) error {
	req := new(validation.CreateNote)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	user, err := getUserFromLocals(c)
	if err != nil {
		return err
	}
	note, err := nc.NoteService.CreateNote(c, req, user.ID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Note created successfully",
		"data":    note,
	})
}

// @Tags         Notes
// @Summary      Get all notes for the user
// @Description  Retrieves a list of all notes belonging to the authenticated user.
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  example.GetNotesResponse
// @Failure      401  {object}  example.Unauthorized
// @Router       /notes [get]
func (nc *NoteController) GetNotes(c *fiber.Ctx) error {
	user, err := getUserFromLocals(c)
	if err != nil {
		return err
	}

	notes, err := nc.NoteService.GetNotesForUser(c, user.ID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"results": len(notes),
		"data":    notes,
	})
}

// @Tags         Notes
// @Summary      Get a single note by ID
// @Description  Retrieves a specific note by its ID. User must be the owner.
// @Security     BearerAuth
// @Produce      json
// @Param        noteId  path      string  true  "Note ID"
// @Success      200     {object}  example.GetNoteResponse
// @Failure      400     {object}  example.InvalidNoteIDError
// @Failure      401     {object}  example.Unauthorized
// @Failure      404     {object}  example.NotFound "Note not found"
// @Router       /notes/{noteId} [get]
func (nc *NoteController) GetNote(c *fiber.Ctx) error {
	noteID, err := uuid.Parse(c.Params("noteId"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Note ID")
	}

	user, err := getUserFromLocals(c)
	if err != nil {
		return err
	}

	note, err := nc.NoteService.GetNoteByID(c, noteID, user.ID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   note,
	})
}

// @Tags         Notes
// @Summary      Update a note
// @Description  Updates a specific note by its ID. User must be the owner.
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        noteId   path      string                 true  "Note ID"
// @Param        request  body      validation.UpdateNote  true  "Request body"
// @Success      200      {object}  example.UpdateNoteResponse
// @Failure      400      {object}  example.InvalidNoteIDError
// @Failure      401      {object}  example.Unauthorized
// @Failure      404      {object}  example.NotFound "Note not found"
// @Router       /notes/{noteId} [patch]
func (nc *NoteController) UpdateNote(c *fiber.Ctx) error {
	noteID, err := uuid.Parse(c.Params("noteId"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Note ID")
	}

	req := new(validation.UpdateNote)
	if err = c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	user, err := getUserFromLocals(c)
	if err != nil {
		return err
	}

	note, err := nc.NoteService.UpdateNote(c, req, noteID, user.ID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Note updated successfully",
		"data":    note,
	})
}

// @Tags         Notes
// @Summary      Delete a note
// @Description  Deletes a specific note by its ID. User must be the owner.
// @Security     BearerAuth
// @Produce      json
// @Param        noteId  path  string  true  "Note ID"
// @Success      204     "No Content"
// @Failure      400     {object}  example.InvalidNoteIDError
// @Failure      401     {object}  example.Unauthorized
// @Failure      404     {object}  example.NotFound "Note not found"
// @Router       /notes/{noteId} [delete]
func (nc *NoteController) DeleteNote(c *fiber.Ctx) error {
	noteID, err := uuid.Parse(c.Params("noteId"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Note ID")
	}

	user, err := getUserFromLocals(c)
	if err != nil {
		return err
	}

	if err = nc.NoteService.DeleteNote(c, noteID, user.ID); err != nil {
		return err
	}

	return c.Status(fiber.StatusNoContent).JSON(nil)
}

// @Tags         Notes
// @Summary      Upload a file to a note
// @Description  Uploads a single file and attaches it to a note. User must be the owner.
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        noteId  path      string  true  "Note ID"
// @Param        file    formData  file    true  "The file to upload"
// @Success      200     {object}  example.UploadFileResponse
// @Failure      400     {object}  example.InvalidNoteIDError
// @Failure      401     {object}  example.Unauthorized
// @Failure      404     {object}  example.NotFound "Note not found"
// @Router       /notes/{noteId}/upload [post]
func (nc *NoteController) UploadFile(c *fiber.Ctx) error {
	noteID, err := uuid.Parse(c.Params("noteId"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Note ID")
	}

	user, err := getUserFromLocals(c)
	if err != nil {
		return err
	}

	note, err := nc.NoteService.UploadFileToNote(c, noteID, user.ID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "File uploaded and attached to note successfully",
		"data":    note,
	})
}

func getUserFromLocals(c *fiber.Ctx) (*model.User, error) {
	user, ok := c.Locals("user").(*model.User)
	if !ok || user == nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "User context not found")
	}
	return user, nil
}
