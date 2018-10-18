package log

import (
	"bytes"
	"crypto/sha1"
	"regexp"
	"strconv"
	"time"
)

type (
	ErrorLine struct {
		Date time.Time
		Type string
		Pid  int
		Ip   string
		Port int
		Log  string
		Hash []byte
	}

	ErrorFilter struct {
		Date string
		Type string
		Pid  string
		Ip   string
		Port string
	}
)


// validate the split line against the filter to check if it should be kept or discarded
func (f *ErrorFilter) Validate(line []string) bool {
	if 6 != len(line) {
		return false
	}
	
	if "" != f.Date {
		re, err := regexp.Compile(f.Date)
		if nil != err {
			return false
		}
		
		if !re.Match([]byte(line[0])) {
			return false
		}
	}

	if f.Type != "" && line[1] != f.Type {
		return false
	}

	if "" != f.Pid && line[2] != f.Pid {
		return false
	}

	if "" != f.Ip && line[3] != f.Ip {
		return false
	}

	if "" != f.Port && line[4] != f.Port {
		return false
	}

	return true
}


func parseErrorLine(line []string) (eline *ErrorLine) {
	if 6 != len(line) {
		return
	}

	eline = &ErrorLine{}
	eline.Date, _ = time.Parse("Mon Jan 2 15:04:05.999999 2006", line[0])
	eline.Type    = line[1]
	eline.Pid, _  = strconv.Atoi(line[2])
	eline.Ip      = line[3]
	eline.Port, _ = strconv.Atoi(line[4])
	eline.Log     = line[5]

	hash := sha1.New()
	eline.Hash = hash.Sum([]byte(line[5]))

	return
}


func genErrorRegex() string {
	var buffer bytes.Buffer

	buffer.WriteString(`^\[([^\]]+)\]\s?`)
	buffer.WriteString(`\[([^\]]+)\]\s?`)
	buffer.WriteString(`\[pid\s([^\]]+)\]\s?`)
	buffer.WriteString(`\[client\s([\d\.])+(\d|)\]\s?`)
	buffer.WriteString(`(.+)$`)

	return buffer.String()
}


func ParseErrorLog(path string, filter *ErrorFilter) []*ErrorLine {
	scanner, file := scanner(path)
	defer file.Close()

	var lines []*ErrorLine

	for scanner.Scan() {
		line := scanner.Text()
		data, err := regexp.Compile(genErrorRegex())
		if nil != err {
			continue
		}

		parts := data.FindStringSubmatch(line)

		if nil != filter || filter.Validate(parts) {

			if eline := parseErrorLine(parts); nil != eline {
				lines = append(lines, eline)
			}

		}
	}

	return lines
}