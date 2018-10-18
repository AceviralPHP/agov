package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"strings"
)

func CreateArchive(zipPath string, files []string) (err error) {
	zipFile, err := os.Create(zipPath)
	if nil != err { return }

	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	for _, file := range files {
		err = func() (err error){
			tmpf, err := os.Open(file)
			if nil != err { return }

			defer tmpf.Close()

			info, err := tmpf.Stat()
			if nil != err { return }

			header, err := zip.FileInfoHeader(info)
			if nil != err { return }

			header.Name   = strings.Replace(file, ROOT + "\\", "", 1)
			header.Method = zip.Deflate

			writer, err := archive.CreateHeader(header)
			if nil != err { return }


			if _, err = io.Copy(writer, tmpf); nil != err {
				return
			}

			return
		}()

		if nil != err {
			return
		}
	}
}