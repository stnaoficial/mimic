package internal

import (
	"fmt"
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

func (d *Dumper) Dump(entryMap EntryMap) lang.AnalisysMap {
	analisysMap := make(lang.AnalisysMap)

	for pathName, entry := range entryMap {
		newAnalisysMap, err := d.analyzer.Analyze(lang.NewBuffer("<pathname>", pathName))

		if err != nil {
			cli.Logln(cli.LogSeverityError, err.Error())
			os.Exit(0)
		}

		maps.Copy(analisysMap, newAnalisysMap)

		if entry.IsDir() {
			d.dumpDirectory(analisysMap, pathName, entry)
		} else {
			d.dumpFile(analisysMap, pathName, entry)
		}
	}

	for _, analisys := range analisysMap {
		switch dp := analisys.(type) {
		case lang.VariableAnalisys:
			fmt.Printf("%s: %s\n", dp.Name, dp.Value)
		}
	}

	fmt.Println()

	return analisysMap
}

func (d *Dumper) dumpDirectory(_ lang.AnalisysMap, dirName string, entry Entry) {
	if d.debug {
		cli.Logf(cli.LogSeverityWarn, "Dumping directory %s ...\n", dirName)
	}
}

func (d *Dumper) dumpFile(analisysMap lang.AnalisysMap, fileName string, entry Entry) {
	before, isCompilable := strings.CutSuffix(fileName, ".mimic")

	if isCompilable {
		fileName = before
	}

	if d.debug {
		cli.Logf(cli.LogSeverityWarn, "Dumping file %s ...\n", fileName)
	}

	if isCompilable {
		newAnalisysMap, err := d.analyzer.Analyze(lang.NewBuffer(fileName, string(entry.Data)))

		if err != nil {
			cli.Logln(cli.LogSeverityError, err.Error())
			os.Exit(0)
		}

		maps.Copy(analisysMap, newAnalisysMap)
	}
}
