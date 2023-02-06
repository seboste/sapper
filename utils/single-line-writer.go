package utils

import (
	"fmt"
	"io"
)

const escape = "\x1b"

type SingleLineWriter struct {
	writer              io.Writer
	hasEndedWithNewline *bool
}

func MakeSingleLineWriter(w io.Writer) SingleLineWriter {
	fmt.Fprintf(w, "%s[s", escape) //save terminal position
	newFalse := false
	return SingleLineWriter{writer: w, hasEndedWithNewline: &newFalse}
}

func (slw SingleLineWriter) Write(p []byte) (int, error) {
	requiredByteCount := len(p)
	cleanupRequired := *slw.hasEndedWithNewline

	//remove tailing new lines
	var startIndex int
	for startIndex = len(p) - 1; startIndex >= 0 && p[startIndex] == '\n'; startIndex-- {
		*slw.hasEndedWithNewline = true
	}
	p = p[:startIndex+1]

	//find last row
	for startIndex = len(p) - 1; startIndex >= 0 && p[startIndex] != '\n'; startIndex-- {
	}
	p = p[startIndex+1:]

	if startIndex != -1 { //newline has been found
		cleanupRequired = true
	}

	if cleanupRequired {
		slw.Cleanup()
	}

	_, err := slw.writer.Write(p)
	return requiredByteCount, err
}

func (slw SingleLineWriter) Cleanup() {
	fmt.Fprintf(slw.writer, "%s[u", escape) //return to saved position
	fmt.Fprintf(slw.writer, "%s[J", escape) //erase everything that has been written
	*slw.hasEndedWithNewline = false
}

var _ io.Writer = SingleLineWriter{}
