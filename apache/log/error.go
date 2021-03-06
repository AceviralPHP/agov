package log

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type (
	ErrorLine struct {
		Date     time.Time
		Type     string
		Pid      int
		Ip       string
		Port     int
		Log      string
		Referrer string
		Hash     string
	}

	ErrorFilter struct {
		Date     string
		Type     string
		Pid      string
		Ip       string
		Port     string
		Referrer string
	}
)


// validate the split line against the filter to check if it should be kept or discarded
func (f *ErrorFilter) Validate(line []string) bool {
	if 7 > len(line) {
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

	if f.Type != "" && line[2] != f.Type {
		return false
	}

	if "" != f.Pid && line[3] != f.Pid {
		return false
	}

	if "" != f.Ip && line[4] != f.Ip {
		return false
	}

	if "" != f.Port && line[5] != f.Port {
		return false
	}

	if "" != f.Referrer && line[7] != f.Referrer {
		return false
	}

	return true
}


func parseErrorLine(line []string) (eline *ErrorLine) {
	if 7 > len(line) {
		return
	}

	if 7 == len(line) {
		line = append(line, "")
	}

	eline = &ErrorLine{}
	eline.Date, _  = time.Parse("Mon Jan 2 15:04:05.999999 2006", line[1])
	eline.Type     = line[2]
	eline.Pid, _   = strconv.Atoi(line[3])
	eline.Ip       = line[4]
	eline.Port, _  = strconv.Atoi(line[5])
	eline.Log      = line[6]
	eline.Referrer = line[7]

	hash := sha1.New()
	hash.Write([]byte(line[6]))
	eline.Hash = fmt.Sprintf("%x", hash.Sum(nil))

	return
}


func genErrorRegex() string {
	var buffer bytes.Buffer

	buffer.WriteString(`^\[([^\]]+)\]\s?`)
	buffer.WriteString(`\[([^\]]+)\]\s?`)
	buffer.WriteString(`\[pid\s([^\]]+)\]\s?`)
	buffer.WriteString(`\[client\s([\d\.]+):(\d+)\]\s?`)
	buffer.WriteString(`(.+?)`)
	buffer.WriteString(`[,\sreferer\:(.+)]?$`)

	return buffer.String()
}


func ParseErrorLog(path string, filter *ErrorFilter) []*ErrorLine {
	scanner, file := scanner(path)
	defer file.Close()

	var lines []*ErrorLine
	regstring := genErrorRegex()
	regCheck, err := regexp.Compile(regstring)
	if nil != err {
		fmt.Println(err)
		return lines
	}

	for scanner.Scan() {
		line  := scanner.Text()
		parts := regCheck.FindStringSubmatch(line)

		if nil == filter || filter.Validate(parts) {

			if eline := parseErrorLine(parts); nil != eline {
				lines = append(lines, eline)
			}

		}
	}

	return lines
}


func ParseErrorString(logText string, filter *ErrorFilter) []*ErrorLine {
	text := strings.Split(logText, "\n")

	var lines []*ErrorLine
	regstring := genErrorRegex()
	regCheck, err := regexp.Compile(regstring)
	if nil != err {
		fmt.Println(err)
		return lines
	}


	for _, line := range text {
		parts := regCheck.FindStringSubmatch(line)

		if nil == filter || filter.Validate(parts) {

			if eline := parseErrorLine(parts); nil != eline {
				lines = append(lines, eline)
			}

		}
	}

	return lines
}