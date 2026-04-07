package entity

import (
	"io"
	"time"
)

type UploadInput struct {
	EntityType  string
	EntityID    string
	FileName    string
	ContentType string
	File        io.Reader
	Size        int64
}

type UploadResult struct {
	Key        string
	PublicURL  string
	Size       int64
	UploadedAt time.Time
}
