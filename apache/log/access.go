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
	if 10 != len(line) {
		return false
	}

	if "" != f.Date {
		re, err := regexp.Compile(f.Date)
		if nil != err {
			return false
		}

		if !re.Match([]byte(line[2])) {
			return false
		}
	}

	if "" != f.Ip && line[1] != f.Ip {
		return false
	}

	if "" != f.Verb && line[3] != f.Verb {
		return false
	}

	if "" != f.Path && line[4] != f.Path {
		return false
	}

	if "" != f.Protocol && line[5] != f.Protocol {
		return false
	}

	if "" != f.Code && line[6] != f.Code {
		return false
	}

	if "" != f.Size && line[7] != f.Size {
		return false
	}

	if "" != f.Referrer && line[8] != f.Referrer {
		return false
	}

	if "" != f.UserAgent && line[9] != f.UserAgent {
		return false
	}

	return true
}


func parseAccessLine(line []string) (aline *AccessLine) {
	if 10 != len(line) {
		return
	}

	aline = &AccessLine{}
	aline.Ip        = line[1]
	aline.Date, _   = time.Parse("Mon Jan 2 15:04:05.999999 2006", line[2])
	aline.Verb      = line[3]
	aline.Path      = line[4]
	aline.Protocol  = line[5]
	aline.Code, _   = strconv.Atoi(line[6])
	aline.Size, _   = strconv.Atoi(line[7])
	aline.Referrer  = line[8]
	aline.UserAgent = line[9]

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

		parts := data.FindStringSubmatch(line)

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