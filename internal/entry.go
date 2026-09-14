package internal

import "os"

type FileSystemEntry struct {
	Name string
	Info os.FileInfo
	Size int
	Data []byte
}

type FileSystemMap = map[string]FileSystemEntry

func NewFileEntry(name string, info os.FileInfo, data []byte) FileSystemEntry {
	return FileSystemEntry{
		Name: name,
		Info: info,
		Size: len(data),
		Data: data,
	}
}

func NewDirectoryEntry(name string, info os.FileInfo) FileSystemEntry {
	return FileSystemEntry{
		Name: name,
		Info: info,
		Size: 0,
		Data: nil,
	}
}

func (e *FileSystemEntry) IsDir() bool {
	return e.Info.IsDir()
}

func (e *FileSystemEntry) IsFile() bool {
	return !e.Info.IsDir()
}
