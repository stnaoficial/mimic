package internal

import (
	"fmt"
	"maps"
	"mimic/internal/cli"
	"mimic/internal/lang"
	"os"
	"strings"
	"text/tabwriter"
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

func (d *Dumper) Dump(entryMap FileSystemMap) (lang.AnalisysMap, error) {
	analisysMap := make(lang.AnalisysMap)

	for pathName, entry := range entryMap {
		newAnalisysMap, err := d.analyzer.Analyze(lang.NewBuffer("<pathname>", pathName))

		if err != nil {
			return nil, err
		}

		maps.Copy(analisysMap, newAnalisysMap)

		if entry.IsDir() {
			d.dumpDirectory(analisysMap, pathName, entry)
		} else if err := d.dumpFile(analisysMap, pathName, entry); err != nil {
			return nil, err
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)

	fmt.Fprintln(w, "NAME\tVALUE")

	for _, analisys := range analisysMap {
		switch dp := analisys.(type) {
		case lang.VariableAnalisys:
			fmt.Fprintf(w, "%s\t%s\n", dp.Name, dp.Value)
		}
	}

	fmt.Fprintln(w)

	w.Flush()

	return analisysMap, nil
}

func (d *Dumper) dumpDirectory(_ lang.AnalisysMap, dirName string, _ FileSystemEntry) {
	// allow debug
	if d.debug {
		cli.Logf(cli.LogSeverityWarn, "Dumping directory %s ...\n", dirName)
	}
}

func (d *Dumper) dumpFile(analisysMap lang.AnalisysMap, fileName string, entry FileSystemEntry) error {
	before, isCompilable := strings.CutSuffix(fileName, ".mimic")

	if isCompilable {
		fileName = before
	}

	// allow debug
	if d.debug {
		cli.Logf(cli.LogSeverityWarn, "Dumping file %s ...\n", fileName)
	}

	if isCompilable {
		newAnalisysMap, err := d.analyzer.Analyze(lang.NewBuffer(fileName, string(entry.Data)))

		if err != nil {
			return err
		}

		maps.Copy(analisysMap, newAnalisysMap)
	}

	return nil
}
