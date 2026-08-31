package internal

import (
	"maps"
	"mimic/internal/cli"
	"mimic/internal/lang"
	"os"
	"strings"
)

type Dumper struct {
	analyzer *lang.Analyzer

	debug bool
}

func NewDumper(analyzer *lang.Analyzer, debug bool) *Dumper {
	return &Dumper{
		analyzer: analyzer,

		debug: debug,
	}
}

func (d *Dumper) Dump(inputEntryMap EntryMap) lang.AnalisysMap {
	outputAnalisysMap := make(lang.AnalisysMap)

	for pathName, entry := range inputEntryMap {
		analisysMap, err := d.analyzer.Analyze(lang.NewBuffer("<pathname>", pathName))

		if err != nil {
			cli.Logln(cli.LogSeverityError, err.Error())
			os.Exit(0)
		}

		maps.Copy(outputAnalisysMap, analisysMap)

		if entry.IsDir() {
			d.dumpDirectory(outputAnalisysMap, pathName, entry)
		} else {
			d.dumpFile(outputAnalisysMap, pathName, entry)
		}
	}

	return outputAnalisysMap
}

func (d *Dumper) dumpDirectory(_ lang.AnalisysMap, dirName string, entry Entry) {
	if d.debug {
		cli.Logf(cli.LogSeverityWarn, "Dumping directory %s ...\n", dirName)
	}
}

func (d *Dumper) dumpFile(outputAnalisysMap lang.AnalisysMap, fileName string, entry Entry) {
	before, isCompilable := strings.CutSuffix(fileName, ".mimic")

	if isCompilable {
		fileName = before
	}

	if d.debug {
		cli.Logf(cli.LogSeverityWarn, "Dumping file %s ...\n", fileName)
	}

	if isCompilable {
		analisysMap, err := d.analyzer.Analyze(lang.NewBuffer(fileName, string(entry.Data)))

		if err != nil {
			cli.Logln(cli.LogSeverityError, err.Error())
			os.Exit(0)
		}

		maps.Copy(outputAnalisysMap, analisysMap)
	}
}
