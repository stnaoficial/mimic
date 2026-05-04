package core

import (
	"fmt"
	"mimic/core/cli"
	"mimic/core/lang"
	"mimic/core/util"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Writer struct {
	comp *lang.Compiler

	entryMap EntryMap
}

func NewWriter(comp *lang.Compiler) *Writer {
	return &Writer{
		comp: comp,

		entryMap: make(EntryMap),
	}
}

func (w *Writer) Write(targetPath string, entryMap EntryMap) EntryMap {
	cli.Log(fmt.Sprintf("Writing files to directory %s ...", targetPath), cli.LogSeverityInfo)

	w.entryMap = make(EntryMap)

	w.defineGlobalVars()

	for pathName, entry := range entryMap {
		w.defineLocalVars(pathName, entry)

		pathName = w.comp.Compile(lang.NewBuffer("<pathname>", pathName))
		pathName = filepath.Join(targetPath, pathName)

		if entry.IsDir() {
			w.writeDirectory(pathName, entry)
		} else {
			w.writeFile(pathName, entry)
		}
	}

	return w.entryMap
}

func (w *Writer) defineGlobalVars() {
	now := time.Now()

	w.comp.Env.Vars["__TIMESTAMP__"] = fmt.Sprintf("%d", now.Unix())

	w.comp.Env.Vars["__DATE__"] = now.Format("2006-01-02")
	w.comp.Env.Vars["__TIME__"] = now.Format("15:04:05")
	w.comp.Env.Vars["__DATETIME__"] = now.Format("2006-01-02T15:04:05Z")

	w.comp.Env.Vars["__YEAR__"] = now.Format("2006")
	w.comp.Env.Vars["__MONTH__"] = now.Format("01")
	w.comp.Env.Vars["__DAY__"] = now.Format("02")

	w.comp.Env.Vars["__HOUR__"] = now.Format("15")
	w.comp.Env.Vars["__MINUTE__"] = now.Format("04")
	w.comp.Env.Vars["__SECOND__"] = now.Format("05")

	ns := now.Nanosecond()

	w.comp.Env.Vars["__MILLISECOND__"] = fmt.Sprintf("%03d", ns/1_000_000)
	w.comp.Env.Vars["__MICROSECOND__"] = fmt.Sprintf("%06d", ns/1_000)
	w.comp.Env.Vars["__NANOSECOND__"] = fmt.Sprintf("%09d", ns)
}

func (w *Writer) defineLocalVars(pathName string, entry Entry) {
	w.comp.Env.Vars["__UID__"] = uuid.NewString()

	w.comp.Env.Vars["__16_DIGIT__"] = util.RandDigit(16)
	w.comp.Env.Vars["__8_DIGIT__"] = util.RandDigit(8)
	w.comp.Env.Vars["__4_DIGIT__"] = util.RandDigit(4)
	w.comp.Env.Vars["__2_DIGIT__"] = util.RandDigit(2)

	w.comp.Env.Vars["__BASEPATH__"] = filepath.Dir(pathName)
	w.comp.Env.Vars["__BASENAME__"] = filepath.Base(pathName)

	delete(w.comp.Env.Vars, "__DIRNAME__")
	delete(w.comp.Env.Vars, "__FILENAME__")
	delete(w.comp.Env.Vars, "__FILEDATA__")

	if entry.IsDir() {
		w.comp.Env.Vars["__DIRNAME__"] = pathName
	} else {
		w.comp.Env.Vars["__FILENAME__"] = pathName
		w.comp.Env.Vars["__FILEDATA__"] = string(entry.Data)
	}
}

func (w *Writer) writeDirectory(dirName string, entry Entry) {
	cli.Log(fmt.Sprintf("Writing directory %s ...", dirName), cli.LogSeverityInfo)

	if err := os.MkdirAll(dirName, 0755); err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to create directory %s", dirName), cli.LogSeverityError)
	}

	w.entryMap[dirName] = entry
}

func (w *Writer) writeFile(fileName string, entry Entry) {
	before, isCompilable := strings.CutSuffix(fileName, ".mimic")

	if isCompilable {
		fileName = before
	}

	cli.Log(fmt.Sprintf("Writing file %s ...", fileName), cli.LogSeverityInfo)

	if isCompilable {
		entry.Data = []byte(w.comp.Compile(lang.NewBuffer(fileName, string(entry.Data))))
	}

	dirName := filepath.Dir(fileName)

	if err := os.MkdirAll(dirName, 0755); err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to create directory %s", dirName), cli.LogSeverityError)
	}

	if err := os.WriteFile(fileName, entry.Data, 0644); err != nil {
		cli.LogAndExit(fmt.Sprintf("Unable to write file %s", fileName), cli.LogSeverityError)
	}

	w.entryMap[fileName] = entry
}
