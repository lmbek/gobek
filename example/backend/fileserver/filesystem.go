package fileserver

import (
	"io"
	"net/http"
	"os"
)

type FileSystem struct {
	FileSystem       http.FileSystem
	ReadDirBatchSize int // configuration parameter for `Readdir` func
}

func (fileSystem FileSystem) Open(name string) (http.File, error) {
	file, err := fileSystem.FileSystem.Open(name)

	if err != nil {
		// File doesnt exist, returning error (probably 404 - os.FileNotFound)
		// Custom 404 can be given here
		return nil, err
	}
	return NeuteredStatFile{File: file, readDirBatchSize: fileSystem.ReadDirBatchSize}, nil
}

type NeuteredStatFile struct {
	http.File
	readDirBatchSize int
}

func (neuteredStatFile NeuteredStatFile) Stat() (os.FileInfo, error) {
	stat, statError := neuteredStatFile.File.Stat()

	if statError != nil {
		return nil, statError
	}

	if stat.IsDir() {
	LOOP:
		for {
			fileList, fileListError := neuteredStatFile.File.Readdir(neuteredStatFile.readDirBatchSize)
			switch fileListError {
			case io.EOF:
				break LOOP
			case nil:
				for _, file := range fileList {
					if file.Name() == "index.html" {
						return stat, fileListError // returning index.html (if there is one in searched directory, and no file was searched for)
					}
				}
			default:
				return nil, fileListError
			}
		}
		return nil, os.ErrNotExist // returning 404 error here
	}

	return stat, statError // returning file found
}
