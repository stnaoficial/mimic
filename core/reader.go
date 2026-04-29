package core

import (
	"fmt"
	"mimic/core/cli"
	"mimic/core/util"
	"os"
	"path/filepath"
)

type Reader struct {
	basepath string
	fileMap  util.FileMap
}

func NewReader() *Reader {
	return &Reader{
		basepath: "",
		fileMap:  make(util.FileMap),
	}
}

func (r *Reader) readDirectory(dirname string) {
	cli.Log(fmt.Sprintf("Reading directory %s ...", dirname), cli.LogSeverityInfo)

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
	cli.Log(fmt.Sprintf("Reading file %s ...", filename), cli.LogSeverityInfo)

	filedata, err := os.ReadFile(filename)

	if err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to obtain data from file %s", filename), cli.LogSeverityError)
	}

	relpath, err := filepath.Rel(r.basepath, filename)

	if err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to obtain relative path for file %s", filename), cli.LogSeverityError)
	}

	r.fileMap[relpath] = string(filedata)
}

func (r *Reader) Read(sourcepath string) util.FileMap {
	for k := range r.fileMap {
		delete(r.fileMap, k)
	}

	var basepath = sourcepath

	pathinfo, err := os.Stat(sourcepath)

	if err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to obtain information about path %s", sourcepath), cli.LogSeverityError)
	}

	if !pathinfo.IsDir() {
		basepath = filepath.Dir(sourcepath)
	}

	r.basepath = basepath

	if pathinfo.IsDir() {
		r.readDirectory(sourcepath)
	} else {
		r.readFile(sourcepath)
	}

	return r.fileMap
}
