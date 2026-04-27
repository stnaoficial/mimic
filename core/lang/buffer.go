package lang

import (
	"fmt"
	"strings"
)

type Buffer struct {
	Name   string
	Data   []rune
	Index  int
	Line   int
	Column int
}

func NewBuffer(name string, data []byte) *Buffer {
	return &Buffer{
		Name:   name,
		Data:   []rune(string(data)),
		Index:  0,
		Line:   1,
		Column: 1,
	}
}

func (l *Buffer) Peek() rune {
	if l.Index >= len(l.Data) {
		return 0
	}

	return l.Data[l.Index]
}

func (l *Buffer) Advance() rune {
	ch := l.Peek()

	l.Index++

	if ch == '\n' {
		l.Line++
		l.Column = 1
	} else {
		l.Column++
	}

	return ch
}

func (f *Buffer) GetLineText() string {
	start := f.Index

	for start > 0 && f.Data[start-1] != '\n' {
		start--
	}

	end := f.Index

	for end < len(f.Data) && f.Data[end] != '\n' {
		end++
	}

	return f.expandLineTabs(f.Data[start:end])
}

func (f *Buffer) expandLineTabs(line []rune) string {
	const tabSize = 4

	var result strings.Builder
	col := 0

	for _, ch := range line {
		if ch == '\t' {
			spaces := tabSize - (col % tabSize)
			result.WriteString(strings.Repeat(" ", spaces))
			col += spaces
		} else {
			result.WriteRune(ch)
			col++
		}
	}

	return result.String()
}

func (f *Buffer) buildCaretLine(line string) string {
	const tabSize = 4

	var result strings.Builder

	currentCol := 1

	for _, ch := range line {
		if currentCol >= f.Column {
			break
		}

		if ch == '\t' {
			spaces := tabSize - ((currentCol - 1) % tabSize)
			result.WriteString(strings.Repeat(" ", spaces))
			currentCol += spaces
			continue
		}

		result.WriteRune(' ')
		currentCol++
	}

	result.WriteRune('^')

	return result.String()
}

func (f *Buffer) String() string {
	lineText := f.GetLineText()

	return fmt.Sprintf(
		"%s:%d:%d\n%s\n%s",
		f.Name, f.Line, f.Column,
		lineText, f.buildCaretLine(lineText))
}
