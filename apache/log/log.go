package log

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
)


func scanner(path string) (*bufio.Scanner, *os.File) {
	path, _ = filepath.Abs(path)
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f)
	return scanner, f
}