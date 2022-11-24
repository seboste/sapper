package core

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

type section struct {
	name      string
	verb      string
	lineBegin int //first line number after the begin tag
	lineEnd   int //the line number of the end tag
	content   string
}

type tag struct {
	name  string
	begin bool
	verb  string
}

func readTag(line string) *tag {
	tagExp := regexp.MustCompile(`<<<SAPPER\s*SECTION\s*(BEGIN|END)(\s*(APPEND|REPLACE))?\s*(.*?)>>>`)

	matches := tagExp.FindStringSubmatch(line)
	if len(matches) != 5 {
		return nil
	}

	t := tag{name: matches[4], verb: matches[3]}

	if matches[1] == "BEGIN" {
		t.begin = true
	}

	return &t
}

func readSections(data string) ([]section, error) {
	sections := []section{}
	scanner := bufio.NewScanner(strings.NewReader(data))
	var currentSection *section = nil
	lineCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		if tag := readTag(line); tag != nil {
			if tag.begin == true {
				if currentSection != nil {
					return sections, fmt.Errorf("found nested begin tag %s in existing section %s in line %d", tag.name, currentSection.name, lineCount)
				}
				currentSection = &section{
					name:      tag.name,
					verb:      tag.verb,
					lineBegin: lineCount + 1,
				}
			} else { //tag.begin == false, i.e. end of a section
				if currentSection == nil {
					return sections, fmt.Errorf("found end tag %s without preceeding begin tag in line %d", tag.name, lineCount)
				}
				if tag.name != currentSection.name {
					return sections, fmt.Errorf("found end tag %s does not match the begin tag %s in line %d", tag.name, currentSection.name, lineCount)
				}
				currentSection.lineEnd = lineCount
				sections = append(sections, *currentSection)
				currentSection = nil
			}
		} else {
			if currentSection != nil {
				currentSection.content = currentSection.content + fmt.Sprintln(line)
			}
		}
		lineCount = lineCount + 1
	}
	return sections, nil
}
