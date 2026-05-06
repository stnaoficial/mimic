package core

type EntryType int

const (
	EntryTypeFile EntryType = iota
	EntryTypeDirectory
)

type Entry struct {
	Name string
	Type EntryType
	Size int
	Data []byte
}

type EntryMap = map[string]Entry

func NewFileEntry(name string, data []byte) Entry {
	return Entry{
		Name: name,
		Type: EntryTypeFile,
		Size: len(data),
		Data: data,
	}
}

func NewDirectoryEntry(name string) Entry {
	return Entry{
		Name: name,
		Type: EntryTypeDirectory,
		Size: 0,
		Data: nil,
	}
}

func (n Entry) IsDir() bool {
	return n.Type == EntryTypeDirectory
}

func (n Entry) IsFile() bool {
	return n.Type == EntryTypeFile
}
