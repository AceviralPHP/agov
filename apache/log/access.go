package log

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type (
	AccessLine struct {
		Ip        string
		Date      time.Time
		Verb      string
		Path      string
		Protocol  string
		Code      int
		Size      int
		Referrer  string
		UserAgent string
	}

	AccessFilter struct {
		Ip        string
		Date      string
		Verb      string
		Path      string
		Protocol  string
		Code      string
		Size      string
		Referrer  string
		UserAgent string
	}
)


// validate the split line against the filter to check if it should be kept or discarded
func (f *AccessFilter) Validate(line []string) bool {
	if 9 != len(line) {
		return false
	}

	if "" != f.Date {
		re, err := regexp.Compile(f.Date)
		if nil != err {
			return false
		}

		if !re.Match([]byte(line[1])) {
			return false
		}
	}

	if "" != f.Ip && line[0] != f.Ip {
		return false
	}

	if "" != f.Verb && line[2] != f.Verb {
		return false
	}

	if "" != f.Path && line[3] != f.Path {
		return false
	}

	if "" != f.Protocol && line[4] != f.Protocol {
		return false
	}

	if "" != f.Code && line[5] != f.Code {
		return false
	}

	if "" != f.Size && line[6] != f.Size {
		return false
	}

	if "" != f.Referrer && line[7] != f.Referrer {
		return false
	}

	if "" != f.UserAgent && line[8] != f.UserAgent {
		return false
	}

	return true
}


func parseAccessLine(line []string) (aline *AccessLine) {
	if 9 != len(line) {
		return
	}

	aline = &AccessLine{}
	aline.Ip        = line[0]
	aline.Date, _   = time.Parse("Mon Jan 2 15:04:05.999999 2006", line[1])
	aline.Verb      = line[2]
	aline.Path      = line[3]
	aline.Protocol  = line[4]
	aline.Code, _   = strconv.Atoi(line[5])
	aline.Size, _   = strconv.Atoi(line[6])
	aline.Referrer  = line[7]
	aline.UserAgent = line[8]

	return
}


func genAccessRegex() string {
	var buffer bytes.Buffer

	buffer.WriteString(`^([\d\.]+)\s`)
	buffer.WriteString(`[^\s]+\s+?`)
	buffer.WriteString(`[^\s]+\s+?`)
	buffer.WriteString(`\[([^\]]+)\]\s?`)
	buffer.WriteString(`"([A-Z]+)\s?`)
	buffer.WriteString(`([^"]+)\s`)
	buffer.WriteString(`([^"]+)"\s?`)
	buffer.WriteString(`(\d+)\s`)
	buffer.WriteString(`(\d+)\s`)
	buffer.WriteString(`"([^"]+)"\s`)
	buffer.WriteString(`"(.+)"$`)

	return buffer.String()
}


func ParseAccessLog(path string, filter *AccessFilter) []*AccessLine {
	scanner, file := scanner(path)
	defer file.Close()

	var lines []*AccessLine

	for scanner.Scan() {
		line := scanner.Text()
		data, err := regexp.Compile(genAccessRegex())
		if nil != err {
			fmt.Println(err)
			continue
		}

		parts := data.FindAllString(line, -1)

		if nil != filter || filter.Validate(parts) {

			if aline := parseAccessLine(parts); nil != aline {
				lines = append(lines, aline)
			} else {
				fmt.Println("Failed to parse")
			}
		} else if !filter.Validate(parts) {
			fmt.Println("failed to validate")
		}

	}

	return lines
}