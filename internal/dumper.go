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
	config *Config
	dumper *lang.Dumper

	dumpMap lang.DumpMap

	debug bool
}

func NewDumper(config *Config, debug bool) *Dumper {
	env := lang.NewEnvironment()
	maps.Copy(env.Vars, config.Variables)
	maps.Copy(env.Prompts, config.Prompts)

	expr := lang.NewExpressionConfigurable(config.ExprOpen, config.ExprClose)

	dumper := lang.NewDumperConfigurable(env, expr)

	return &Dumper{
		config: config,
		dumper: dumper,

		dumpMap: make(lang.DumpMap),

		debug: debug,
	}
}

func (d *Dumper) Dump(entryMap EntryMap) lang.DumpMap {
	d.dumpMap = make(lang.DumpMap)

	for _, targetPath := range d.config.TargetPath.Values {
		if d.debug {
			cli.Logf(cli.LogSeverityWarn, "Dumping files for directory %s ...\n", targetPath)
		}

		for pathName, entry := range entryMap {
			dumpMap, err := d.dumper.Dump(lang.NewBuffer("<pathname>", pathName))

			if err != nil {
				cli.Logln(cli.LogSeverityError, err.Error())
				os.Exit(0)
			}

			maps.Copy(d.dumpMap, dumpMap)

			if entry.IsDir() {
				d.dumpDirectory(pathName, entry)
			} else {
				d.dumpFile(pathName, entry)
			}
		}
	}

	for _, dump := range d.dumpMap {
		switch dp := dump.(type) {
		case lang.VariableDump:
			fmt.Printf("%s: %s\n", dp.Name, dp.Value)
		}
	}

	return d.dumpMap
}

func (d *Dumper) dumpDirectory(dirName string, entry Entry) {
	if d.debug {
		cli.Logf(cli.LogSeverityWarn, "Dumping directory %s ...\n", dirName)
	}
}

func (d *Dumper) dumpFile(fileName string, entry Entry) {
	before, isCompilable := strings.CutSuffix(fileName, ".mimic")

	if isCompilable {
		fileName = before
	}

	if d.debug {
		cli.Logf(cli.LogSeverityWarn, "Dumping file %s ...\n", fileName)
	}

	if isCompilable {
		dumpMap, err := d.dumper.Dump(lang.NewBuffer(fileName, string(entry.Data)))

		if err != nil {
			cli.Logln(cli.LogSeverityError, err.Error())
			os.Exit(0)
		}

		maps.Copy(d.dumpMap, dumpMap)
	}
}
