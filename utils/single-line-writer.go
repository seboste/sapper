package utils

import (
	"io"
	"unicode/utf8"
)

type SingleLineWriter struct {
	writer              io.Writer
	hasEndedWithNewline *bool
	noOfRunesWritten    *int
}

func MakeSingleLineWriter(w io.Writer) SingleLineWriter {
	newFalse := false
	newZeroInt := 0
	return SingleLineWriter{writer: w, hasEndedWithNewline: &newFalse, noOfRunesWritten: &newZeroInt}
}

func (slw SingleLineWriter) Write(p []byte) (int, error) {
	requiredByteCount := len(p)
	err := error(nil)
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
		_, err = slw.Cleanup()
		if err != nil {
			return 0, err
		}
	}

	_, err = slw.writer.Write(p)
	if err != nil {
		return 0, err
	}

	*slw.noOfRunesWritten = *slw.noOfRunesWritten + utf8.RuneCountInString(string(p))
	return requiredByteCount, err
}

func (slw SingleLineWriter) Cleanup() (int, error) {
	p := make([]byte, *slw.noOfRunesWritten+1)
	for i := range p {
		p[i] = '\b'
	}
	p[len(p)-1] = '\r' //make sure to return to the beginning of the line
	*slw.hasEndedWithNewline = false
	*slw.noOfRunesWritten = 0
	return slw.writer.Write(p)
}

var _ io.Writer = SingleLineWriter{}
