package log

import (
	"bufio"
	"log"
	"os"

)


func scanner(path string) (*bufio.Scanner, *os.File) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f)
	return scanner, f
}