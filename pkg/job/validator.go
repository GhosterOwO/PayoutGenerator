package job

import (
	"bytes"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"future.net.co/rngcollector/pkg/parser"
)

// Remove line which container element 0.
type ValidatorJob struct {
	File  string
	Chunk *bytes.Buffer
	info  *parser.Service
	mux   *sync.Mutex
}

func NewValidatorJob(filename string, b *bytes.Buffer, s *parser.Service) Job {
	c := bufferPool.Get()
	chunk := c.(*bytes.Buffer)
	chunk.ReadFrom(b)
	return &ValidatorJob{
		mux:   aquireMutex(filename),
		File:  filename,
		Chunk: chunk,
		info:  s,
	}
}

func (j *ValidatorJob) Execute() {
	c := bufferPool.Get()
	chunk := c.(*bytes.Buffer)
	defer chunk.Reset()
	defer bufferPool.Put(chunk)
	defer j.Chunk.Reset()
	defer bufferPool.Put(j.Chunk)
	for {
		line, err := j.Chunk.ReadString('\n')
		s := strings.Trim(line, "\n\r")
		a := strings.Split(strings.TrimSpace(s), " ")
		if len(a) == int(j.info.Selections()) {
			for len(a) > 0 {
				v, err := strconv.ParseInt(a[len(a)-1], 10, 64)
				if err != nil {
					break
				} else if v > j.info.RangeMax() {
					break
				} else if v < j.info.RangeMin() {
					break
				}
				a = a[:len(a)-1]
			}
			if len(a) == 0 {
				chunk.WriteString(line)
			}
		}
		if err == io.EOF {
			break
		} else if err != nil {
			log.Print(err)
			break
		}
	}
	if chunk.Len() <= 0 {
		return
	}
	j.mux.Lock()
	defer j.mux.Unlock()
	const flags = os.O_APPEND | os.O_WRONLY | os.O_CREATE
	if out, err := os.OpenFile(j.File, flags, 0644); err == nil {
		out.WriteString(chunk.String())
		defer out.Close()
	}
}
