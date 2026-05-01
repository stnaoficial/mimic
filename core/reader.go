package core

import (
	"fmt"
	"mimic/core/cli"
	"mimic/core/util"
	"os"
	"path/filepath"
)

type Reader struct {
	fileMap util.FileMap
}

func NewReader() *Reader {
	return &Reader{
		fileMap: make(util.FileMap),
	}
}

func (r *Reader) readDirectory(basePath string, dirName string) {
	cli.Log(fmt.Sprintf("Reading directory %s ...", dirName), cli.LogSeverityInfo)

	fileNames, err := util.DirectoryWalk(dirName)

	if err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to walk into directory %s", dirName), cli.LogSeverityError)
	}

	if len(fileNames) == 0 {
		cli.LogAndExit(fmt.Sprintf("No .mimic files found in directory %s", dirName), cli.LogSeverityWarn)
	}

	for _, fileName := range fileNames {
		r.readFile(basePath, fileName)
	}
}

func (r *Reader) readFile(basePath string, fileName string) {
	cli.Log(fmt.Sprintf("Reading file %s ...", fileName), cli.LogSeverityInfo)

	fileData, err := os.ReadFile(fileName)

	if err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to obtain data from file %s", fileName), cli.LogSeverityError)
	}

	relPath, err := filepath.Rel(basePath, fileName)

	if err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to obtain relative path for file %s", fileName), cli.LogSeverityError)
	}

	r.fileMap[relPath] = string(fileData)
}

func (r *Reader) Read(sourcePath string) util.FileMap {
	for k := range r.fileMap {
		delete(r.fileMap, k)
	}

	sourceInfo, err := os.Stat(sourcePath)

	if err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to obtain information about path %s", sourcePath), cli.LogSeverityError)
	}

	if sourceInfo.IsDir() {
		r.readDirectory(sourcePath, sourcePath)
	} else {
		r.readFile(filepath.Dir(sourcePath), sourcePath)
	}

	return r.fileMap
}
