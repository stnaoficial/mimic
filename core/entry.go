package core

import "os"

type Entry struct {
	Name string
	Info os.FileInfo
	Size int
	Data []byte
}

type EntryMap = map[string]Entry

func NewFileEntry(name string, info os.FileInfo, data []byte) Entry {
	return Entry{
		Name: name,
		Info: info,
		Size: len(data),
		Data: data,
	}
}

func NewDirectoryEntry(name string, info os.FileInfo) Entry {
	return Entry{
		Name: name,
		Info: info,
		Size: 0,
		Data: nil,
	}
}

func (n Entry) IsDir() bool {
	return n.Info.IsDir()
}

func (n Entry) IsFile() bool {
	return !n.Info.IsDir()
}
