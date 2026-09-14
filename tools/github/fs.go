package github

import (
	"os"
	"time"
)

type RemoteEntry struct {
	Path string
	Info os.FileInfo
}

type RemoteFileInfo struct {
	name    string
	size    int64
	modTime time.Time
	mode    os.FileMode
	isDir   bool
}

func (f *RemoteFileInfo) Name() string       { return f.name }
func (f *RemoteFileInfo) Size() int64        { return f.size }
func (f *RemoteFileInfo) Mode() os.FileMode  { return f.mode }
func (f *RemoteFileInfo) ModTime() time.Time { return f.modTime }
func (f *RemoteFileInfo) IsDir() bool        { return f.isDir }
func (f *RemoteFileInfo) Sys() any           { return nil }

type RemoteEntryMap map[string]RemoteEntry
