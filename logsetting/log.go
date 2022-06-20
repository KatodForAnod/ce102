package logsetting

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"time"
)

const dirName = "logs"

func LogInit() error {
	err := os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		log.Fatal(err)
		return err
	}

	fileLogName := time.Now().Format("2006-01-02") + ".txt"
	f, err := os.OpenFile(dirName+"/"+fileLogName, os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(f)
	return nil
}

func OpenLastLogFile() (*os.File, error) {
	fileInfo, err := ioutil.ReadDir(dirName)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	sort.SliceStable(fileInfo, func(i, j int) bool {
		return fileInfo[i].ModTime().Unix() > fileInfo[j].ModTime().Unix()
	})

	if len(fileInfo) == 0 {
		err := errors.New("not found log files")
		log.Println(err)
		return nil, err
	}

	file, err := os.Open(dirName + "/" + fileInfo[0].Name())
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return file, nil
}

func GetNLastLines(file *os.File, nRows int) (lines []string, err error) {
	defer file.Close()

	var cursor int64 = 0
	stat, _ := file.Stat()
	filesize := stat.Size()
	for i := 0; i < nRows; i++ {
		line := ""
		for {
			cursor -= 1
			file.Seek(cursor, io.SeekEnd)

			char := make([]byte, 1)
			file.Read(char)
			if cursor != -1 && (char[0] == 10 || char[0] == 13) { // stop if we find a line
				break
			}

			line = fmt.Sprintf("%s%s", string(char), line) // there is more efficient way
			if cursor == -filesize {                       // stop if we are at the begining
				break
			}
		}

		lines = append(lines, line)
		if cursor == -filesize {
			return lines, nil
		}
	}

	return
}
