package core

type EntryType int

const (
	EntryTypeFile EntryType = iota
	EntryTypeDirectory
)

type Entry struct {
	Name string
	Type EntryType
	Data []byte
}

type EntryMap = map[string]Entry

func NewFileEntry(name string, data []byte) Entry {
	return Entry{
		Name: name,
		Type: EntryTypeFile,
		Data: data,
	}
}

func NewDirectoryEntry(name string) Entry {
	return Entry{
		Name: name,
		Type: EntryTypeDirectory,
		Data: nil,
	}
}

func (n Entry) IsDir() bool {
	return n.Type == EntryTypeDirectory
}

func (n Entry) IsFile() bool {
	return n.Type == EntryTypeFile
}
