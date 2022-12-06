package job

import (
	"bytes"
	"io"
	"log"
	"os"
)

type LineCountJob struct {
	File   string
	result chan int64
}

func NewLineCountJob(filename string, r chan int64) Job {
	return &LineCountJob{
		File:   filename,
		result: r,
	}
}

func (j *LineCountJob) Execute() {
	f, err := os.Open(j.File)
	if err != nil {
		return
	}
	defer f.Close()
	count := int64(0)
	buf := make([]byte, 32*1024)
	lineSep := []byte{'\n'}
	for err == nil {
		var c int
		c, err = f.Read(buf)
		count += int64(bytes.Count(buf[:c], lineSep))
	}
	if err != io.EOF {
		log.Print(err)
		return
	}
	j.result <- count
}
