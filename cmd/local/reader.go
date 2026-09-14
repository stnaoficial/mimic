package local

import (
	"fmt"
	"mimic/internal/cli"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/tabwriter"
)

type Reader struct {
	entryPathsRead []string

	debug bool
}

func NewReader(debug bool) *Reader {
	return &Reader{
		entryPathsRead: []string{},

		debug: debug,
	}
}

func (r *Reader) List(sourcePaths []string) error {
	r.entryPathsRead = []string{}

	for _, sourcePath := range sourcePaths {
		// allow debug
		if r.debug {
			cli.Logf(cli.LogSeverityWarn, "Listing source path %s ...\n", sourcePath)
		}

		sourceInfo, err := os.Stat(sourcePath)

		if err != nil {
			return err
		}

		if sourceInfo.IsDir() {
			if err := r.listSourceDirectory(sourcePath); err != nil {
				return err
			}
		} else {
			r.listSourceFile(filepath.Dir(sourcePath))
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)

	fmt.Fprintln(w, "NAME\tPATH")

	templateNames := []string{}

	for _, entryPath := range r.entryPathsRead {
		parts := strings.Split(entryPath, string(os.PathSeparator))

		if len(parts) == 0 {
			continue
		}

		templateName := parts[0]

		if slices.Contains(templateNames, templateName) {
			continue
		}

		fmt.Fprintf(w, "%s\t%s\n", templateName, entryPath)

		templateNames = append(templateNames, templateName)
	}

	fmt.Fprintln(w)

	w.Flush()

	return nil
}

func (r *Reader) listSourceDirectory(dirName string) error {
	// allow debug
	if r.debug {
		cli.Logf(cli.LogSeverityWarn, "Listing directory %s ...\n", dirName)
	}

	entriesFound, err := os.ReadDir(dirName)

	if err != nil {
		return err
	}

	for _, entry := range entriesFound {
		r.listEntry(entry.Name())
	}

	return nil
}

func (r *Reader) listSourceFile(fileName string) {
	// allow debug
	if r.debug {
		cli.Logf(cli.LogSeverityWarn, "Listing file %s ...\n", fileName)
	}

	r.listEntry(fileName)
}

func (r *Reader) listEntry(fileName string) {
	r.entryPathsRead = append(r.entryPathsRead, fileName)
}
