package job

import (
	"bytes"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

// Remove line which include duplicate element.
type NonReplacementJob struct {
	File  string
	Chunk *bytes.Buffer
	mux   *sync.Mutex
}

func NewNonReplacementJob(filename string, b *bytes.Buffer) Job {
	c := bufferPool.Get()
	chunk := c.(*bytes.Buffer)
	chunk.ReadFrom(b)
	return &NonReplacementJob{
		mux:   aquireMutex(filename),
		File:  filename,
		Chunk: chunk,
	}
}

func unique(arr []string) bool {
	occurred := map[string]bool{}
	for e := range arr {
		if occurred[arr[e]] != true {
			occurred[arr[e]] = true
		} else {
			return false
		}
	}
	return true
}

func (j *NonReplacementJob) Execute() {
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
		if unique(a) {
			chunk.WriteString(line)
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
