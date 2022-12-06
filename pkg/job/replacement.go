package job

import (
	"bytes"
	"os"
	"regexp"
	"sync"
)

// Remove line which container element 0.
type ReplacementJob struct {
	File  string
	Chunk *bytes.Buffer
	mux   *sync.Mutex
}

func NewReplacementJob(filename string, b *bytes.Buffer) Job {
	c := bufferPool.Get()
	chunk := c.(*bytes.Buffer)
	chunk.ReadFrom(b)
	return &ReplacementJob{
		mux:   aquireMutex(filename),
		File:  filename,
		Chunk: chunk,
	}
}

func (j *ReplacementJob) Execute() {
	pattern := "[[:space:]]0[[:space:]]"
	re := regexp.MustCompile("(?m)^.*" + pattern + ".*$[\r\n]+")
	chunk := re.ReplaceAllString(j.Chunk.String(), "") + "\n"
	j.mux.Lock()
	defer j.mux.Unlock()
	defer j.Chunk.Reset()
	defer bufferPool.Put(j.Chunk)
	const flags = os.O_APPEND | os.O_WRONLY | os.O_CREATE
	if out, err := os.OpenFile(j.File, flags, 0644); err == nil {
		out.WriteString(chunk)
		defer out.Close()
	}
}
