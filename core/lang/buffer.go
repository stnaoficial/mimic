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

func NewBuffer(name string, data string) *Buffer {
	return &Buffer{
		Name:   name,
		Data:   []rune(data),
		Index:  0,
		Line:   1,
		Column: 1,
	}
}

func (b *Buffer) Peek() rune {
	if b.Index >= len(b.Data) {
		return 0
	}

	return b.Data[b.Index]
}

func (b *Buffer) Advance() rune {
	ch := b.Peek()

	b.Index++

	if ch == '\n' {
		b.Line++
		b.Column = 1
	} else {
		b.Column++
	}

	return ch
}

func (b *Buffer) lineText() string {
	start := b.Index

	for start > 0 && b.Data[start-1] != '\n' {
		start--
	}

	end := b.Index

	for end < len(b.Data) && b.Data[end] != '\n' {
		end++
	}

	return b.expandLineTabs(b.Data[start:end])
}

func (b *Buffer) expandLineTabs(line []rune) string {
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

func (b *Buffer) buildCaretLine(line string) string {
	const tabSize = 4

	var result strings.Builder

	currentCol := 1

	for _, ch := range line {
		if currentCol >= b.Column {
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

func (b *Buffer) String() string {
	lineText := b.lineText()

	return fmt.Sprintf(
		"%s:%d:%d\n%s\n%s",
		b.Name, b.Line, b.Column,
		lineText, b.buildCaretLine(lineText))
}
