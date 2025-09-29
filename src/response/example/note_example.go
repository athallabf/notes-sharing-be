package example

import (
	"time"

	"github.com/google/uuid"
)

// Note represents a single note object for Swagger examples.
type Note struct {
	ID        uuid.UUID `json:"id" example:"a1b2c3d4-e5f6-a7b8-c9d0-e1f2a3b4c5d6"`
	Title     string    `json:"title" example:"My First Note"`
	Content   string    `json:"content" example:"This is the content of the note."`
	FilePath  string    `json:"file_path,omitempty" example:"abcdef-12345.pdf"`
	FileURL   string    `json:"file_url,omitempty" example:"https://s3.tebi.io/your-bucket/abcdef-12345.pdf?..."`
	CreatedAt time.Time `json:"created_at" example:"2025-09-29T10:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2025-09-29T10:00:00Z"`
}

// CreateNoteResponse is the example for a successful note creation.
type CreateNoteResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Note created successfully"`
	Data    Note   `json:"data"`
}

// GetNoteResponse is the example for fetching a single note.
type GetNoteResponse struct {
	Status string `json:"status" example:"success"`
	Data   Note   `json:"data"`
}

// GetNotesResponse is the example for fetching a list of notes.
type GetNotesResponse struct {
	Status  string `json:"status" example:"success"`
	Results int    `json:"results" example:"1"`
	Data    []Note `json:"data"`
}

// UpdateNoteResponse is the example for a successful note update.
type UpdateNoteResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Note updated successfully"`
	Data    Note   `json:"data"`
}

// UploadFileResponse is the example for a successful file upload.
type UploadFileResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"File uploaded and attached to note successfully"`
	Data    Note   `json:"data"`
}

// InvalidNoteIDError is an example for a bad request due to an invalid UUID.
type InvalidNoteIDError struct {
	Code    int    `json:"code" example:"400"`
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Invalid Note ID"`
}
