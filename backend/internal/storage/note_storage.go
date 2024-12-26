package storage

import "cloud_technologies/internal/models"

type NoteStorage interface {
	GetNote(string) (*models.Note, error)
	UploadNote(*models.Note) error
}
