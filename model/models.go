package model

import (
	"os"
	"time"
)

type File struct {
	Path string
	FileInfo FileInfo
}

type FileInfo struct {
	Name    string
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
	IsDir   bool
}


type Stats struct {
	TotalFiles uint64
	MaxFile    struct {
		Size int64
		Path string
	}
	AvgFileSize  float64
	Extensions   []string
	TopExtension string
}

