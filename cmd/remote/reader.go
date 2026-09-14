package remote

import (
	"fmt"
	"mimic/internal/cli"
	"mimic/tools/github"
	"os"
	"text/tabwriter"
)

type Reader struct {
	debug bool
}

func NewReader(debug bool) *Reader {
	return &Reader{
		debug: debug,
	}
}

func (r *Reader) List(sourcePaths []string) error {
	entries := make(map[string]github.ApiEntry)

	for _, sourcePath := range sourcePaths {
		// allow debug
		if r.debug {
			cli.Logf(cli.LogSeverityWarn, "Listing source path %s ...\n", sourcePath)
		}

		entriesFound, err := GithubApiClient.FetchApiContent(sourcePath)

		if err != nil {
			return err
		}

		for _, entry := range entriesFound {
			entries[entry.Name] = entry
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)

	fmt.Fprintln(w, "NAME\tPATH\tURL")

	for entryName, entry := range entries {
		fmt.Fprintf(w, "%s\t%s\t%s\n", entryName, entry.Path, GithubApiClient.ParsePublicEntryUrl(entry.Path))
	}

	fmt.Fprintln(w)

	w.Flush()

	return nil
}
