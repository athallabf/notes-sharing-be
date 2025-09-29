package service

import (
	"app/src/model"
	"app/src/utils"
	"app/src/validation"
	"errors"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type NoteService interface {
	CreateNote(c *fiber.Ctx, req *validation.CreateNote, userID uuid.UUID) (*model.Note, error)
	GetNoteByID(c *fiber.Ctx, noteID, userID uuid.UUID) (*model.Note, error)
	GetNotesForUser(c *fiber.Ctx, userID uuid.UUID) ([]model.Note, error)
	UpdateNote(c *fiber.Ctx, req *validation.UpdateNote, noteID, userID uuid.UUID) (*model.Note, error)
	DeleteNote(c *fiber.Ctx, noteID, userID uuid.UUID) error
	UploadFileToNote(c *fiber.Ctx, noteID, userID uuid.UUID) (*model.Note, error)
}

type noteService struct {
	Log      *logrus.Logger
	DB       *gorm.DB
	Validate *validator.Validate
	S3       S3Service
}

func NewNoteService(db *gorm.DB, validate *validator.Validate, s3 S3Service) NoteService {
	return &noteService{
		Log:      utils.Log,
		DB:       db,
		Validate: validate,
		S3:       s3,
	}
}

func (s *noteService) CreateNote(c *fiber.Ctx, req *validation.CreateNote, userID uuid.UUID) (*model.Note, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	note := &model.Note{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID,
	}

	if err := s.DB.WithContext(c.Context()).Create(note).Error; err != nil {
		s.Log.Errorf("Failed to create note: %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	return note, nil
}

func (s *noteService) GetNoteByID(c *fiber.Ctx, noteID, userID uuid.UUID) (*model.Note, error) {
	var note model.Note
	if err := s.DB.WithContext(c.Context()).Where("id = ? AND user_id = ?", noteID, userID).
		First(&note).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "Note not found")
		}
		s.Log.Errorf("Failed to get note by ID: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if note.FilePath != "" {
		url, err := s.S3.GeneratePresignedURL(note.FilePath)
		if err != nil {
			s.Log.Warnf("Failed to generate presigned URL for note %s: %v", note.ID, err)
		} else {
			note.FileURL = url
		}
	}
	return &note, nil
}

func (s *noteService) GetNotesForUser(c *fiber.Ctx, userID uuid.UUID) ([]model.Note, error) {
	var notes []model.Note
	if err := s.DB.WithContext(c.Context()).
		Order("created_at desc").
		Where("user_id = ?", userID).
		Find(&notes).Error; err != nil {
		s.Log.Errorf("Failed to get notes for user: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	for i := range notes {
		if notes[i].FilePath != "" {
			url, err := s.S3.GeneratePresignedURL(notes[i].FilePath)
			if err != nil {
				s.Log.Warnf("Failed to generate presigned URL for note %s: %v", notes[i].ID, err)
			} else {
				notes[i].FileURL = url
			}
		}
	}

	return notes, nil
}

func (s *noteService) UpdateNote(
	c *fiber.Ctx,
	req *validation.UpdateNote,
	noteID, userID uuid.UUID,
) (*model.Note, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	note, err := s.GetNoteByID(c, noteID, userID)
	if err != nil {
		return nil, err
	}

	note.Title = req.Title
	note.Content = req.Content

	if err = s.DB.WithContext(c.Context()).Save(note).Error; err != nil {
		s.Log.Errorf("Failed to update note: %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	return note, nil
}

func (s *noteService) DeleteNote(c *fiber.Ctx, noteID, userID uuid.UUID) error {
	result := s.DB.WithContext(c.Context()).Where("id = ? AND user_id = ?", noteID, userID).Delete(&model.Note{})
	if result.Error != nil {
		s.Log.Errorf("Failed to delete note: %+v", result.Error)
		return fiber.ErrInternalServerError
	}
	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Note not found")
	}
	return nil
}

func (s *noteService) UploadFileToNote(c *fiber.Ctx, noteID, userID uuid.UUID) (*model.Note, error) {
	note, err := s.GetNoteByID(c, noteID, userID)
	if err != nil {
		return nil, err
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "File not provided")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Unable to open file")
	}
	defer file.Close()

	uniqueFileName := uuid.New().String() + filepath.Ext(fileHeader.Filename)

	if err = s.S3.UploadFile(file, uniqueFileName); err != nil {
		s.Log.Errorf("Failed to upload file to S3: %v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Unable to save file")
	}

	note.FilePath = uniqueFileName
	if err = s.DB.WithContext(c.Context()).Save(note).Error; err != nil {
		s.Log.Errorf("Failed to update note with file path: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return note, nil
}
