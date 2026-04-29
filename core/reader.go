package core

import (
	"fmt"
	"mimic/core/cli"
	"mimic/core/util"
	"os"
)

type Reader struct {
	fileMap util.FileMap
}

func NewReader() *Reader {
	return &Reader{
		fileMap: make(util.FileMap),
	}
}

func (r *Reader) readDirectory(dirname string) {
	cli.Log(fmt.Sprintf("Reading directory %s...", dirname), cli.LogSeverityInfo)

	filenames, err := util.DirectoryWalk(dirname)

	if err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to walk into directory %s", dirname), cli.LogSeverityError)
	}

	if len(filenames) == 0 {
		cli.LogAndExit(fmt.Sprintf("No .mimic files found in directory %s", dirname), cli.LogSeverityWarn)
	}

	for _, filename := range filenames {
		r.readFile(filename)
	}
}

func (r *Reader) readFile(filename string) {
	cli.Log(fmt.Sprintf("Reading file %s...", filename), cli.LogSeverityInfo)

	filedata, err := os.ReadFile(filename)

	if err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to obtain data from file %s", filename), cli.LogSeverityError)
	}

	r.fileMap[filename] = string(filedata)
}

func (r *Reader) Read(pathname string, pathInfo os.FileInfo) util.FileMap {
	for k := range r.fileMap {
		delete(r.fileMap, k)
	}

	if pathInfo.IsDir() {
		r.readDirectory(pathname)
	} else {
		r.readFile(pathname)
	}

	return r.fileMap
}
