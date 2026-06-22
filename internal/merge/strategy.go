package merge

import (
	"strings"
)

type Strategy struct{}

func NewStrategy() *Strategy {
	return &Strategy{}
}

func (s *Strategy) Merge(old, new string) (string, error) {
	oldLines := splitLines(old)
	newLines := splitLines(new)

	// find the longest common prefix (the outer header)
	prefixLen := 0

	for prefixLen < len(oldLines) && prefixLen < len(newLines) && oldLines[prefixLen] == newLines[prefixLen] {
		prefixLen++
	}

	// find the longest common suffix (the outer footer)
	suffixLen := 0
	minIndent := -1

	for suffixLen < len(oldLines)-prefixLen && suffixLen < len(newLines)-prefixLen {
		oldIdx := len(oldLines) - 1 - suffixLen
		newIdx := len(newLines) - 1 - suffixLen

		if oldLines[oldIdx] != newLines[newIdx] {
			break
		}

		line := oldLines[oldIdx]

		if strings.TrimSpace(line) != "" {
			indent := countIndent(line)
			if minIndent == -1 {
				minIndent = indent
			} else if indent > minIndent {
				// indentation increased (moving upwards), meaning we hit an inner scope block.
				break
			} else {
				minIndent = indent
			}
		}

		suffixLen++
	}

	// extract the components
	oldPrefix := oldLines[:prefixLen]
	oldSuffix := oldLines[len(oldLines)-suffixLen:]

	oldMid := oldLines[prefixLen : len(oldLines)-suffixLen]
	newMid := newLines[prefixLen : len(newLines)-suffixLen]

	// assemble the merged buffer sequence
	var resultLines []string

	resultLines = append(resultLines, oldPrefix...)
	resultLines = append(resultLines, oldMid...)

	// ensure a structural blank line separates the two unique blocks if they both contain content
	if len(oldMid) > 0 && len(newMid) > 0 {
		if strings.TrimSpace(oldMid[len(oldMid)-1]) != "" && strings.TrimSpace(newMid[0]) != "" {
			resultLines = append(resultLines, "")
		}
	}

	resultLines = append(resultLines, newMid...)
	resultLines = append(resultLines, oldSuffix...)

	return strings.Join(resultLines, "\n"), nil
}

func splitLines(s string) []string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return strings.Split(s, "\n")
}

func countIndent(line string) int {
	count := 0

	for _, ch := range line {
		if ch == ' ' {
			count++
		} else if ch == '\t' {
			count += 4
		} else {
			break
		}
	}
	return count
}
